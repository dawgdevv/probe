import { useState } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { Link } from 'react-router-dom';
import toast from 'react-hot-toast';
import { listProjects, createProject, type Project } from '../api/client';

export default function Dashboard() {
  const queryClient = useQueryClient();
  const [showForm, setShowForm] = useState(false);
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');

  const { data: projects, isLoading } = useQuery({
    queryKey: ['projects'],
    queryFn: listProjects,
  });

  const mutation = useMutation({
    mutationFn: () => createProject(name, description),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['projects'] });
      setShowForm(false);
      setName('');
      setDescription('');
      toast.success('Project created!');
    },
    onError: () => toast.error('Failed to create project'),
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim()) return;
    mutation.mutate();
  };

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold">Projects</h2>
        <button
          onClick={() => setShowForm(!showForm)}
          className="bg-[var(--accent)] hover:bg-[var(--accent-hover)] text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors cursor-pointer border-none"
        >
          + New Project
        </button>
      </div>

      {showForm && (
        <form
          onSubmit={handleSubmit}
          className="bg-[var(--bg-secondary)] border border-[var(--border)] rounded-lg p-4 mb-6"
        >
          <div className="mb-3">
            <label className="block text-sm text-[var(--text-secondary)] mb-1">
              Project Name
            </label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="My API Tests"
              className="w-full bg-[var(--bg-tertiary)] border border-[var(--border)] rounded-md px-3 py-2 text-[var(--text-primary)] text-sm outline-none focus:border-[var(--accent)]"
              autoFocus
            />
          </div>
          <div className="mb-3">
            <label className="block text-sm text-[var(--text-secondary)] mb-1">
              Description
            </label>
            <input
              type="text"
              value={description}
              onChange={(e) => setDescription(e.target.value)}
              placeholder="Optional description"
              className="w-full bg-[var(--bg-tertiary)] border border-[var(--border)] rounded-md px-3 py-2 text-[var(--text-primary)] text-sm outline-none focus:border-[var(--accent)]"
            />
          </div>
          <div className="flex gap-2">
            <button
              type="submit"
              disabled={mutation.isPending}
              className="bg-[var(--accent)] hover:bg-[var(--accent-hover)] text-white px-4 py-2 rounded-md text-sm font-medium cursor-pointer border-none transition-colors disabled:opacity-50"
            >
              {mutation.isPending ? 'Creating...' : 'Create'}
            </button>
            <button
              type="button"
              onClick={() => setShowForm(false)}
              className="bg-[var(--bg-tertiary)] text-[var(--text-secondary)] px-4 py-2 rounded-md text-sm cursor-pointer border border-[var(--border)] transition-colors hover:text-[var(--text-primary)]"
            >
              Cancel
            </button>
          </div>
        </form>
      )}

      {isLoading ? (
        <div className="text-[var(--text-secondary)] text-center py-12">Loading...</div>
      ) : !projects || projects.length === 0 ? (
        <div className="text-center py-16">
          <p className="text-5xl mb-4">📁</p>
          <p className="text-[var(--text-secondary)] text-lg">No projects yet</p>
          <p className="text-[var(--text-secondary)] text-sm">
            Create a project to start organizing your API tests
          </p>
        </div>
      ) : (
        <div className="grid gap-4 grid-cols-1 md:grid-cols-2 lg:grid-cols-3">
          {projects.map((project: Project) => (
            <Link
              key={project.id}
              to={`/projects/${project.id}`}
              className="bg-[var(--bg-secondary)] border border-[var(--border)] rounded-lg p-5 no-underline transition-all hover:border-[var(--accent)] hover:shadow-lg hover:shadow-blue-500/10"
            >
              <h3 className="text-[var(--text-primary)] text-lg font-semibold mb-1">
                {project.name}
              </h3>
              <p className="text-[var(--text-secondary)] text-sm m-0">
                {project.description || 'No description'}
              </p>
              <p className="text-[var(--text-secondary)] text-xs mt-3 m-0">
                Created {new Date(project.created_at).toLocaleDateString()}
              </p>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}
