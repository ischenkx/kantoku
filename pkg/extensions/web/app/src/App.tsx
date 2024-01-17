import {Refine} from "@refinedev/core";
import {ThemedLayoutV2, ErrorComponent, RefineThemes, ThemedTitleV2} from "@refinedev/antd";
import routerBindings, {NavigateToResource, UnsavedChangesNotifier} from "@refinedev/react-router-v6";
import dataProvider from "@refinedev/simple-rest";
import {BrowserRouter, Routes, Route, Outlet} from "react-router-dom";
import {AntdInferencer} from "@refinedev/inferencer/antd";
import MemoryOutlinedIcon from '@mui/icons-material/MemoryOutlined';
import StorageOutlinedIcon from '@mui/icons-material/StorageOutlined';
import "@refinedev/antd/dist/reset.css";
import {TaskShow} from "./entities/task/taskShow";
import {TaskList} from "./entities/task/taskList";
import { Header } from "./components/header";
import {ColorModeContextProvider} from "./contexts/color-mode";
import {ResourceList} from "./entities/resource/resourceList";
import {ResourceShow} from "./entities/resource/resourceShow";
import {ResourceCreate} from "./entities/resource/resourceCreate";
import {TaskCreate} from "./entities/task/taskCreate";
import {Sandbox} from "./entities/sandbox/sandbox";

const App: React.FC = () => {
    return (
        <BrowserRouter>
            <ColorModeContextProvider>
                <Refine
                    routerProvider={routerBindings}
                    // dataProvider={dataProvider("http://127.0.1.1:3030")}
                    dataProvider={dataProvider("http://127.0.0.1:3030")}
                    // dataProvider={dataProvider("http://localhost:3000")}
                    resources={[
                        {
                            name: "tasks",
                            list: "/tasks",
                            show: "/tasks/show/:id",
                            create: "/tasks/create",
                            meta: {
                                icon: <MemoryOutlinedIcon/>,
                            }
                        },
                        {
                            name: "resources",
                            list: "/resources",
                            show: "/resources/show/:id",
                            create: "/resources/create",
                            meta: {
                                icon: <StorageOutlinedIcon/>
                            }
                        },
                        {
                            name: "sandbox",
                            list: "/sandbox",
                        },
                    ]}
                    options={{
                        syncWithLocation: true,
                        warnWhenUnsavedChanges: true,
                    }}
                >
                    <Routes>
                        <Route
                            element={
                                <ThemedLayoutV2
                                    Header={Header}
                                    Title={({ collapsed }) => (
                                        <ThemedTitleV2
                                            // collapsed is a boolean value that indicates whether the <Sidebar> is collapsed or not
                                            collapsed={collapsed}
                                            text="Kantoku"
                                        />
                                    )}
                                >
                                    <Outlet/>
                                </ThemedLayoutV2>
                            }
                        >
                            <Route index element={<NavigateToResource resource="tasks"/>}/>
                            <Route index element={<NavigateToResource resource="resources"/>}/>
                            <Route path="tasks">
                                <Route index element={<TaskList/>}/>
                                <Route path="show/:id" element={<TaskShow/>}/>
                                <Route path="create" element={<TaskCreate/>}/>
                            </Route>
                            <Route path="resources">
                                <Route index element={<ResourceList/>}/>
                                <Route path="show/:id" element={<ResourceShow/>}/>
                                <Route path="create" element={<ResourceCreate/>}/>
                            </Route>
                            <Route path="sandbox">
                                <Route path=":context" index element={<Sandbox/>}/>
                            </Route>
                            <Route path="*" element={<ErrorComponent/>}/>
                        </Route>
                    </Routes>
                    <UnsavedChangesNotifier/>

                </Refine>
            </ColorModeContextProvider>
        </BrowserRouter>
    );
};

export default App;