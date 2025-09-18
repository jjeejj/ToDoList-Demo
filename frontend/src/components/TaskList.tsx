'use client';

import React, { useState, useEffect } from 'react';
import { Trash2, CheckCircle, Circle } from 'lucide-react';
import { todoClient } from '../lib/client';
import { Task, GetTasksRequest, DeleteTaskRequest, UpdateTaskRequest } from '../gen/todolist/v1/todolist_pb';

interface TaskListProps {
  refreshTrigger: number;
}

export default function TaskList({ refreshTrigger }: TaskListProps) {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(false);

  const fetchTasks = async () => {
    try {
      setLoading(true);
      const request = new GetTasksRequest({});
      const response = await todoClient.getTasks(request);
      setTasks(response.tasks || []);
    } catch (error) {
      console.error('Failed to fetch tasks:', error);
    } finally {
      setLoading(false);
    }
  };

  const deleteTask = async (id: string) => {
    try {
      const request = new DeleteTaskRequest({ id });
      const response = await todoClient.deleteTask(request);
      if (response.success) {
        setTasks(tasks.filter(task => task.id !== id));
      }
    } catch (error) {
      console.error('Error deleting task:', error);
    }
  };

  const toggleTaskCompletion = async (id: string, completed: boolean) => {
    try {
      const request = new UpdateTaskRequest({ id, completed: !completed });
      const response = await todoClient.updateTask(request);
      if (response.success && response.task) {
        setTasks(tasks.map(task => 
          task.id === id ? response.task! : task
        ));
      }
    } catch (error) {
      console.error('Error updating task:', error);
    }
  };

  useEffect(() => {
    fetchTasks();
  }, [refreshTrigger]);

  if (loading) {
    return (
      <div className="flex justify-center items-center py-8">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>
    );
  }



  if (tasks.length === 0) {
    return (
      <div className="text-center py-8 text-gray-500">
        <Circle className="mx-auto h-12 w-12 mb-4 opacity-50" />
        <p>No tasks yet. Add your first task above!</p>
      </div>
    );
  }

  return (
    <div className="space-y-3">
      {tasks.map((task) => (
        <div
          key={task.id}
          className={`bg-white border border-gray-200 rounded-lg p-4 shadow-sm hover:shadow-md transition-shadow ${
            task.completed ? 'opacity-75' : ''
          }`}
        >
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <button
                onClick={() => toggleTaskCompletion(task.id, task.completed)}
                className="flex-shrink-0 transition-colors"
              >
                {task.completed ? (
                  <CheckCircle className="h-5 w-5 text-green-500" />
                ) : (
                  <Circle className="h-5 w-5 text-gray-400 hover:text-green-500" />
                )}
              </button>
              <div>
                <h3 className={`font-medium ${
                  task.completed 
                    ? 'text-gray-500 line-through' 
                    : 'text-gray-900'
                }`}>
                  {task.text}
                </h3>
              </div>
            </div>
            <button
              onClick={() => deleteTask(task.id)}
              className="text-red-500 hover:text-red-700 p-2 rounded-full hover:bg-red-50 transition-colors"
              title="Delete task"
            >
              <Trash2 className="h-4 w-4" />
            </button>
          </div>
        </div>
      ))}
    </div>
  );
}