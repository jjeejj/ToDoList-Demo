import { createConnectTransport } from "@connectrpc/connect-web";
import { createClient } from "@connectrpc/connect";
import { TodoService } from "../gen/todolist/v1/todolist_connect";

// Create the transport for connecting to the backend
const transport = createConnectTransport({
  baseUrl: "http://localhost:8080", // Backend server URL
  useBinaryFormat: true, // Use binary protobuf format for better performance
});

// Create the client for TodoService
export const todoClient = createClient(TodoService, transport);

// Export types for use in components
export type { Task, AddTaskRequest, GetTasksResponse, DeleteTaskRequest } from "../gen/todolist/v1/todolist_pb";