'use client';

import React, { useState } from 'react';
import { Plus } from 'lucide-react';
import { todoClient } from '../lib/client';
import { AddTaskRequest } from '../gen/todolist/v1/todolist_pb';

interface AddTaskProps {
  onTaskAdded: () => void;
}

export default function AddTask({ onTaskAdded }: AddTaskProps) {
  const [text, setText] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!text.trim()) return;

    try {
      setLoading(true);
      const request = new AddTaskRequest({ text: text.trim() });
      const response = await todoClient.addTask(request);
      if (response.success) {
        setText('');
        onTaskAdded();
      }
    } catch (error) {
      console.error('Failed to add task:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="bg-white rounded-lg shadow-md p-6 mb-6">
      <h2 className="text-xl font-semibold text-gray-800 mb-4">Add New Task</h2>
      
      <form onSubmit={handleSubmit} className="space-y-4">
        <div>
          <label htmlFor="text" className="block text-sm font-medium text-gray-700 mb-1">
            Task Text *
          </label>
          <input
            type="text"
            id="text"
            value={text}
            onChange={(e) => setText(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
            placeholder="Enter task text..."
            disabled={loading}
          />
        </div>
        
        <button
          type="submit"
          disabled={loading || !text.trim()}
          className="w-full bg-blue-600 text-white py-2 px-4 rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
        >
          <Plus className="w-4 h-4" />
          {loading ? "Adding..." : "Add Task"}
        </button>
      </form>
    </div>
  );
}