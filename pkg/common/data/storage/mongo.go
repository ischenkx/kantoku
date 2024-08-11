package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log/slog"
)

var (
	readCommands = []string{
		"aggregate",
		"count",
		"distinct",
		"geoSearch",
		"find",
		"getMore",
	}
	writeCommands = []string{
		"delete",
		"findAndModify",
		"insert",
		"update",
	}
)

type MongoStorage struct {
	Collection *mongo.Collection
	Logger     *slog.Logger
}

func (s *MongoStorage) Settings(ctx context.Context) (Settings, error) {
	return Settings{
		Type: "mongo",
		Meta: make(map[string]any),
	}, nil
}

// Exec runs commands as described at https://www.mongodb.com/docs/manual/reference/command/
//
// All commands are collection local, so you can't make cross-collection aggregation
// TODO: might be a good idea to allow cross-collection requests
func (s *MongoStorage) Exec(ctx context.Context, command Command) ([]Document, error) {
	switch command.Operation {
	case "aggregate":
		command.Params = append([]Param{{"aggregate", s.Collection.Name()}}, command.Params...)
	case "count":
		command.Params = append([]Param{{"count", s.Collection.Name()}}, command.Params...)
	case "distinct":
		command.Params = append([]Param{{"distinct", s.Collection.Name()}}, command.Params...)
	case "geoSearch":
		command.Params = append([]Param{{"geoSearch", s.Collection.Name()}}, command.Params...)
	case "delete":
		command.Params = append([]Param{{"delete", s.Collection.Name()}}, command.Params...)
	case "find":
		command.Params = append([]Param{{"find", s.Collection.Name()}}, command.Params...)
	case "findAndModify":
		command.Params = append([]Param{{"aggregate", s.Collection.Name()}}, command.Params...)
	case "getMore":
		command.Params = append([]Param{{"getMore", s.Collection.Name()}}, command.Params...)
	case "insert":
		command.Params = append([]Param{{"insert", s.Collection.Name()}}, command.Params...)
	case "update":
		command.Params = append([]Param{{"update", s.Collection.Name()}}, command.Params...)
	default:
		return nil, fmt.Errorf("unknown operation '%s'", command.Operation)
	}

	if rawIsReadOnly, ok := command.Meta["readonly"]; ok {
		isReadOnly, ok := rawIsReadOnly.(bool)
		if !ok {
			return nil, fmt.Errorf("readonly mode must be set to 'true' or 'false'")
		}

		if isReadOnly {
			if !lo.Contains(readCommands, command.Operation) {
				return nil, fmt.Errorf("operation '%s' is not allowed (it violates the read-only constraint)", command.Operation)
			}
		}
	}

	runCommand := bson.D(lo.Map(command.Params, func(param Param, _ int) bson.E {
		return bson.E{param.Name, param.Value}
	}))

	if s.Logger != nil {
		marshalledCommand, err := json.Marshal(runCommand)
		if err != nil {
			s.Logger.Debug("failed to marshal the command",
				slog.String("error", err.Error()))
		} else {
			s.Logger.Debug("exec",
				slog.String("cmd", string(marshalledCommand)))
		}
	}

	result := s.Collection.Database().RunCommand(ctx, runCommand)
	if result.Err() != nil {
		return nil, result.Err()
	}

	var doc Document
	if err := result.Decode(&doc); err != nil {
		return nil, fmt.Errorf("error getting document '%s': %v", command.Operation, err)
	}

	return []Document{doc}, nil
}
