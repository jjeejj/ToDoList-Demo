"use client";

import { useState } from "react";
import AddTask from "../components/AddTask";
import TaskList from "../components/TaskList";
import { ListTodo } from "lucide-react";

export default function Home() {
  const [refreshTrigger, setRefreshTrigger] = useState(0);

  const handleTaskAdded = () => {
    // Trigger a refresh of the task list
    setRefreshTrigger(prev => prev + 1);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 py-8">
      <div className="max-w-4xl mx-auto px-4">
        {/* Header */}
        <div className="text-center mb-8">
          <div className="flex items-center justify-center space-x-3 mb-4">
            <ListTodo className="h-8 w-8 text-blue-600" />
            <h1 className="text-4xl font-bold text-gray-900">
              Todo List
            </h1>
          </div>
          <p className="text-gray-600 text-lg">
            Organize your tasks efficiently with our modern todo application
          </p>
        </div>

        {/* Main Content */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Add Task Section */}
          <div className="lg:col-span-1">
            <AddTask onTaskAdded={handleTaskAdded} />
          </div>
          
          {/* Task List Section */}
          <div className="lg:col-span-2">
            <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
              <h2 className="text-xl font-semibold text-gray-900 mb-6">
                Your Tasks
              </h2>
              <TaskList refreshTrigger={refreshTrigger} />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
