import { useState } from 'react';
import { useParams, Link } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { getProject, listSuites, createSuite, type TestSuite } from '../api/client';

export default function ProjectView() {
  const { projectId } = useParams<{ projectId: string }>();
  const id = Number(projectId);
  const queryClient = useQueryClient();
  const [showForm, setShowForm] = useState(false);
  const [name, setName] = useState('');
  const [yamlContent, setYamlContent] = useState(DEFAULT_YAML);

  const { data: project } = useQuery({
    queryKey: ['project', id],
    queryFn: () => getProject(id),
  });

  const { data: suites, isLoading } = useQuery({
    queryKey: ['suites', id],
    queryFn: () => listSuites(id),
  });

  const mutation = useMutation({
    mutationFn: () => createSuite(id, name, yamlContent),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['suites', id] });
      setShowForm(false);
      setName('');
      setYamlContent(DEFAULT_YAML);
      toast.success('Test suite created!');
    },
    onError: (err: any) =>
      toast.error(err?.response?.data?.error || 'Failed to create suite'),
  });

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!name.trim() || !yamlContent.trim()) return;
    mutation.mutate();
  };

  return (
    <div>
      {/* Breadcrumb */}
      <div className="flex items-center gap-2 text-sm text-[var(--text-secondary)] mb-4">
        <Link to="/" className="hover:text-[var(--text-primary)] no-underline text-[var(--text-secondary)]">
          Projects
        </Link>
        <span>/</span>
        <span className="text-[var(--text-primary)]">{project?.name ?? '...'}</span>
      </div>

      <div className="flex items-center justify-between mb-6">
        <div>
          <h2 className="text-2xl font-bold mb-1">{project?.name}</h2>
          {project?.description && (
            <p className="text-[var(--text-secondary)] text-sm m-0">{project.description}</p>
          )}
        </div>
        <button
          onClick={() => setShowForm(!showForm)}
          className="bg-[var(--accent)] hover:bg-[var(--accent-hover)] text-white px-4 py-2 rounded-lg text-sm font-medium transition-colors cursor-pointer border-none"
        >
          + New Suite
        </button>
      </div>

      {showForm && (
        <form
          onSubmit={handleSubmit}
          className="bg-[var(--bg-secondary)] border border-[var(--border)] rounded-lg p-4 mb-6"
        >
          <div className="mb-3">
            <label className="block text-sm text-[var(--text-secondary)] mb-1">
              Suite Name
            </label>
            <input
              type="text"
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="e.g. User API Tests"
              className="w-full bg-[var(--bg-tertiary)] border border-[var(--border)] rounded-md px-3 py-2 text-[var(--text-primary)] text-sm outline-none focus:border-[var(--accent)]"
              autoFocus
            />
          </div>
          <div className="mb-3">
            <label className="block text-sm text-[var(--text-secondary)] mb-1">
              YAML Test Definition
            </label>
            <textarea
              value={yamlContent}
              onChange={(e) => setYamlContent(e.target.value)}
              rows={16}
              className="w-full bg-[var(--bg-tertiary)] border border-[var(--border)] rounded-md px-3 py-2 text-[var(--text-primary)] text-sm font-mono outline-none focus:border-[var(--accent)] resize-y"
              spellCheck={false}
            />
          </div>
          <div className="flex gap-2">
            <button
              type="submit"
              disabled={mutation.isPending}
              className="bg-[var(--accent)] hover:bg-[var(--accent-hover)] text-white px-4 py-2 rounded-md text-sm font-medium cursor-pointer border-none transition-colors disabled:opacity-50"
            >
              {mutation.isPending ? 'Creating...' : 'Create Suite'}
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
      ) : !suites || suites.length === 0 ? (
        <div className="text-center py-16">
          <p className="text-5xl mb-4">🧪</p>
          <p className="text-[var(--text-secondary)] text-lg">No test suites yet</p>
          <p className="text-[var(--text-secondary)] text-sm">
            Create a test suite with your YAML test definitions
          </p>
        </div>
      ) : (
        <div className="grid gap-4 grid-cols-1 md:grid-cols-2">
          {suites.map((suite: TestSuite) => (
            <Link
              key={suite.id}
              to={`/suites/${suite.id}`}
              className="bg-[var(--bg-secondary)] border border-[var(--border)] rounded-lg p-5 no-underline transition-all hover:border-[var(--accent)] hover:shadow-lg hover:shadow-blue-500/10"
            >
              <h3 className="text-[var(--text-primary)] text-lg font-semibold mb-2">
                {suite.name}
              </h3>
              <p className="text-[var(--text-secondary)] text-xs m-0">
                Updated {new Date(suite.updated_at).toLocaleDateString()}
              </p>
            </Link>
          ))}
        </div>
      )}
    </div>
  );
}

const DEFAULT_YAML = `env:
  base_url: https://jsonplaceholder.typicode.com

tests:
  - name: Get all posts
    request:
      method: GET
      path: /posts
    expect:
      status: 200

  - name: Get single post
    request:
      method: GET
      path: /posts/1
    expect:
      status: 200
`;
