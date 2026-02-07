import { useParams, Link } from 'react-router-dom';
import { useQuery } from '@tanstack/react-query';
import { getTestRun, getTestResults, type TestResult } from '../api/client';

export default function RunView() {
  const { runId } = useParams<{ runId: string }>();
  const id = Number(runId);

  const { data: run } = useQuery({
    queryKey: ['run', id],
    queryFn: () => getTestRun(id),
  });

  const { data: results } = useQuery({
    queryKey: ['runResults', id],
    queryFn: () => getTestResults(id),
  });

  return (
    <div>
      {/* Breadcrumb */}
      <div className="flex items-center gap-2 text-sm text-[var(--text-secondary)] mb-4">
        <Link to="/" className="hover:text-[var(--text-primary)] no-underline text-[var(--text-secondary)]">
          Projects
        </Link>
        <span>/</span>
        {run && (
          <>
            <Link
              to={`/suites/${run.suite_id}`}
              className="hover:text-[var(--text-primary)] no-underline text-[var(--text-secondary)]"
            >
              Suite
            </Link>
            <span>/</span>
          </>
        )}
        <span className="text-[var(--text-primary)]">Run #{runId}</span>
      </div>

      {run && (
        <>
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-2xl font-bold">Test Run #{run.id}</h2>
            <span
              className={`text-sm font-semibold px-3 py-1.5 rounded-lg ${
                run.status === 'passed'
                  ? 'bg-green-500/20 text-[var(--success)]'
                  : run.status === 'failed'
                  ? 'bg-red-500/20 text-[var(--danger)]'
                  : 'bg-yellow-500/20 text-[var(--warning)]'
              }`}
            >
              {run.status.toUpperCase()}
            </span>
          </div>

          {/* Summary cards */}
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
            <InfoCard label="Total Tests" value={run.total_tests} />
            <InfoCard label="Passed" value={run.passed_tests} color="var(--success)" />
            <InfoCard label="Failed" value={run.failed_tests} color="var(--danger)" />
            <InfoCard
              label="Started"
              value={new Date(run.started_at).toLocaleTimeString()}
              small
            />
          </div>
        </>
      )}

      {/* Results table */}
      <h3 className="text-lg font-semibold mb-3">Test Results</h3>
      {!results || results.length === 0 ? (
        <p className="text-[var(--text-secondary)] text-sm">No results available.</p>
      ) : (
        <div className="bg-[var(--bg-secondary)] border border-[var(--border)] rounded-lg overflow-hidden">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-[var(--border)]">
                <th className="text-left px-4 py-3 text-[var(--text-secondary)] font-medium">
                  Status
                </th>
                <th className="text-left px-4 py-3 text-[var(--text-secondary)] font-medium">
                  Test Name
                </th>
                <th className="text-left px-4 py-3 text-[var(--text-secondary)] font-medium">
                  HTTP Status
                </th>
                <th className="text-left px-4 py-3 text-[var(--text-secondary)] font-medium">
                  Duration
                </th>
                <th className="text-left px-4 py-3 text-[var(--text-secondary)] font-medium">
                  Error
                </th>
              </tr>
            </thead>
            <tbody>
              {results.map((result: TestResult) => (
                <tr
                  key={result.id}
                  className="border-b border-[var(--border)] last:border-none"
                >
                  <td className="px-4 py-3">
                    <span className="text-lg">{result.passed ? '✔' : '✖'}</span>
                  </td>
                  <td className="px-4 py-3 font-medium">{result.test_name}</td>
                  <td className="px-4 py-3">
                    {result.status_code > 0 && (
                      <span className="bg-[var(--bg-tertiary)] px-2 py-0.5 rounded text-xs">
                        {result.status_code}
                      </span>
                    )}
                  </td>
                  <td className="px-4 py-3 text-[var(--text-secondary)]">
                    {result.duration_ms}ms
                  </td>
                  <td className="px-4 py-3 text-[var(--danger)] text-xs max-w-[300px] truncate">
                    {result.error_message}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}

function InfoCard({
  label,
  value,
  color,
  small,
}: {
  label: string;
  value: string | number;
  color?: string;
  small?: boolean;
}) {
  return (
    <div className="bg-[var(--bg-secondary)] border border-[var(--border)] rounded-lg px-4 py-3 text-center">
      <div
        className={`font-bold ${small ? 'text-sm' : 'text-2xl'}`}
        style={{ color: color ?? 'var(--text-primary)' }}
      >
        {value}
      </div>
      <div className="text-xs text-[var(--text-secondary)] mt-1">{label}</div>
    </div>
  );
}
