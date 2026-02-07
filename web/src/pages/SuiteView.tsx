import { useState } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import toast from 'react-hot-toast';
import { getSuite, runSuite, listRuns, type TestRun, type RunResponse } from '../api/client';

export default function SuiteView() {
  const { suiteId } = useParams<{ suiteId: string }>();
  const id = Number(suiteId);
  const navigate = useNavigate();
  const queryClient = useQueryClient();
  const [runResult, setRunResult] = useState<RunResponse | null>(null);

  const { data: suite } = useQuery({
    queryKey: ['suite', id],
    queryFn: () => getSuite(id),
  });

  const { data: runs } = useQuery({
    queryKey: ['runs', id],
    queryFn: () => listRuns(id),
  });

  const runMutation = useMutation({
    mutationFn: () => runSuite(id),
    onSuccess: (data) => {
      setRunResult(data);
      queryClient.invalidateQueries({ queryKey: ['runs', id] });
      toast.success(`Tests ${data.status}: ${data.passed_tests}/${data.total_tests} passed`);
    },
    onError: () => toast.error('Failed to run tests'),
  });

  return (
    <div>
      {/* Breadcrumb */}
      <div className="flex items-center gap-2 text-sm text-[var(--text-secondary)] mb-4">
        <Link to="/" className="hover:text-[var(--text-primary)] no-underline text-[var(--text-secondary)]">
          Projects
        </Link>
        <span>/</span>
        {suite && (
          <>
            <Link
              to={`/projects/${suite.project_id}`}
              className="hover:text-[var(--text-primary)] no-underline text-[var(--text-secondary)]"
            >
              Project
            </Link>
            <span>/</span>
          </>
        )}
        <span className="text-[var(--text-primary)]">{suite?.name ?? '...'}</span>
      </div>

      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold">{suite?.name}</h2>
        <button
          onClick={() => runMutation.mutate()}
          disabled={runMutation.isPending}
          className="bg-[var(--success)] hover:brightness-110 text-white px-5 py-2.5 rounded-lg text-sm font-semibold transition-all cursor-pointer border-none disabled:opacity-50"
        >
          {runMutation.isPending ? '⏳ Running...' : '▶ Run Tests'}
        </button>
      </div>

      {/* YAML Content */}
      {suite && (
        <div className="bg-[var(--bg-secondary)] border border-[var(--border)] rounded-lg mb-6 overflow-hidden">
          <div className="px-4 py-2 border-b border-[var(--border)] text-sm text-[var(--text-secondary)] font-medium">
            Test Definition (YAML)
          </div>
          <pre className="p-4 m-0 text-sm font-mono text-[var(--text-primary)] overflow-x-auto whitespace-pre-wrap">
            {suite.yaml_content}
          </pre>
        </div>
      )}

      {/* Live Run Results */}
      {runResult && (
        <div className="mb-6">
          <h3 className="text-lg font-semibold mb-3">Latest Run Results</h3>
          <div className="flex gap-4 mb-4">
            <StatusBadge
              label="Total"
              value={runResult.total_tests}
              color="var(--accent)"
            />
            <StatusBadge
              label="Passed"
              value={runResult.passed_tests}
              color="var(--success)"
            />
            <StatusBadge
              label="Failed"
              value={runResult.failed_tests}
              color="var(--danger)"
            />
          </div>
          <div className="space-y-2">
            {runResult.results.map((r, i) => (
              <div
                key={i}
                className={`flex items-center justify-between bg-[var(--bg-secondary)] border rounded-lg px-4 py-3 ${
                  r.Passed ? 'border-green-500/30' : 'border-red-500/30'
                }`}
              >
                <div className="flex items-center gap-3">
                  <span className="text-lg">{r.Passed ? '✔' : '✖'}</span>
                  <span className="text-sm font-medium">{r.Name}</span>
                </div>
                <div className="flex items-center gap-4 text-sm text-[var(--text-secondary)]">
                  {r.StatusCode > 0 && (
                    <span className="bg-[var(--bg-tertiary)] px-2 py-0.5 rounded text-xs">
                      {r.StatusCode}
                    </span>
                  )}
                  <span>{(r.Duration / 1_000_000).toFixed(0)}ms</span>
                  {r.Error && (
                    <span className="text-[var(--danger)] text-xs max-w-[200px] truncate">
                      {r.Error}
                    </span>
                  )}
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Run History */}
      <div>
        <h3 className="text-lg font-semibold mb-3">Run History</h3>
        {!runs || runs.length === 0 ? (
          <p className="text-[var(--text-secondary)] text-sm">
            No runs yet. Click "Run Tests" to execute this suite.
          </p>
        ) : (
          <div className="space-y-2">
            {runs.map((run: TestRun) => (
              <div
                key={run.id}
                onClick={() => navigate(`/runs/${run.id}`)}
                className={`flex items-center justify-between bg-[var(--bg-secondary)] border rounded-lg px-4 py-3 cursor-pointer transition-colors hover:border-[var(--accent)] ${
                  run.status === 'passed'
                    ? 'border-green-500/30'
                    : run.status === 'failed'
                    ? 'border-red-500/30'
                    : 'border-[var(--border)]'
                }`}
              >
                <div className="flex items-center gap-3">
                  <span
                    className={`text-xs font-semibold px-2 py-1 rounded ${
                      run.status === 'passed'
                        ? 'bg-green-500/20 text-[var(--success)]'
                        : run.status === 'failed'
                        ? 'bg-red-500/20 text-[var(--danger)]'
                        : 'bg-yellow-500/20 text-[var(--warning)]'
                    }`}
                  >
                    {run.status.toUpperCase()}
                  </span>
                  <span className="text-sm">
                    {run.passed_tests}/{run.total_tests} passed
                  </span>
                </div>
                <span className="text-xs text-[var(--text-secondary)]">
                  {new Date(run.started_at).toLocaleString()}
                </span>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}

function StatusBadge({ label, value, color }: { label: string; value: number; color: string }) {
  return (
    <div
      className="bg-[var(--bg-secondary)] border border-[var(--border)] rounded-lg px-4 py-3 text-center min-w-[80px]"
    >
      <div className="text-2xl font-bold" style={{ color }}>
        {value}
      </div>
      <div className="text-xs text-[var(--text-secondary)] mt-1">{label}</div>
    </div>
  );
}
