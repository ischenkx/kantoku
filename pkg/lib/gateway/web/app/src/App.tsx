import {Refine} from "@refinedev/core";
import {ErrorComponent, RefineThemes} from "@refinedev/antd";
import routerBindings, {NavigateToResource, UnsavedChangesNotifier} from "@refinedev/react-router-v6";
import dataProvider from "@refinedev/simple-rest";
import {BrowserRouter, Routes, Route, Outlet} from "react-router-dom";
import {AntdInferencer} from "@refinedev/inferencer/antd";
import MemoryOutlinedIcon from '@mui/icons-material/MemoryOutlined';
import StorageOutlinedIcon from '@mui/icons-material/StorageOutlined';
import DescriptionOutlinedIcon from '@mui/icons-material/DescriptionOutlined';
import AccountTreeOutlinedIcon from '@mui/icons-material/AccountTreeOutlined';
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
import Flow from "./entities/sandbox/flow";
import {ThemedLayoutV2} from "./components/layout";
import {ThemedTitleV2} from "./components/layout/title";
import {ProviderRouter} from "./providers";

const App: React.FC = () => {
    return (
        <BrowserRouter>
            <ColorModeContextProvider>
                <Refine
                    routerProvider={routerBindings}
                    dataProvider={ProviderRouter}
                    resources={[
                        {
                            name: "tasks",
                            list: "/tasks",
                            show: "/tasks/show/:id",
                            create: "/tasks/create",
                            meta: {
                                icon: <MemoryOutlinedIcon/>,
                            },
                        },
                        {
                            name: "specifications",
                            list: "/tasks/specifications",
                            show: "/tasks/specifications/:id",
                            meta: {
                                icon: <DescriptionOutlinedIcon/>,
                            },
                        },
                        {
                            name: "types",
                            list: "/tasks/specifications/types",
                            show: "/tasks/specifications/types/:id",
                            meta: {
                                icon: <AccountTreeOutlinedIcon/>,
                            },
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
                        {
                            name: "flow",
                            list: "/flow",
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
                            <Route path="specifications">
                                <Route index element={<ResourceList/>}/>
                                <Route path="show/:id" element={<ResourceShow/>}/>
                                <Route path="create" element={<ResourceCreate/>}/>
                            </Route>
                            <Route path="sandbox">
                                <Route path=":context" index element={<Sandbox/>}/>
                            </Route>
                            <Route path="flow">
                                <Route index element={<Flow/>}/>
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