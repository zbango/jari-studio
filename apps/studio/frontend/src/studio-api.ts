import {
  CreateProject,
  CreateRevision,
  CompileRevision,
  EvaluateSchedule,
  ExecuteCard,
  StartTerminal,
  TerminalEvents,
  WriteTerminalInput,
  InterruptTerminal,
  CloseTerminal,
  DefaultWorkspaceDirectory,
  SelectWorkspaceDirectory,
  AvailableAdapters,
  Health,
  ListProjects,
  ListRevisions,
} from "../wailsjs/go/main/App";

export type Project = {
  id: string;
  name: string;
  brief?: string;
  revision: number;
  createdAt: unknown;
  updatedAt: unknown;
};

export type StudioHealth = Record<string, string>;

export type Revision = {
  id: string;
  projectId: string;
  number: number;
  status: string;
  snapshot: string;
  createdAt: unknown;
};

export type CardSpec = {
  id: string;
  goal: string;
  scopePaths: string[];
  forbiddenPaths: string[];
  dependencies?: string[];
  acceptanceCriteria: string[];
  verificationCommands: string[];
  requiredAdapterCapabilities: string[];
  state: string;
};

export type CompiledPlan = {
  revisionId: string;
  requirements: Array<{ id: string; statement: string }>;
  acceptanceCriteria: Array<{ id: string; statement: string; mandatory: boolean; verificationType: string }>;
  cards: CardSpec[];
  blockingGaps?: Array<{ id: string; severity: string; status: string }>;
};

export type ScheduleResult = {
  adapter: { id: string; version: string; transport: string; capabilities: string[] };
  decisions: Array<{ cardId: string; state: string; reasons?: string[] }>;
};

export type AdapterManifest = {
  id: string;
  version: string;
  transport: string;
  capabilities: string[];
};

export type RunResult = {
  cardId: string;
  worktree: string;
  state: string;
  events: Array<{ type: string; data: string }>;
  changedFiles: string[];
};

export type TerminalInfo = { id: string; command: string; directory: string; state: string };
export type TerminalEvent = { type: string; data: string };

export function isNativeRuntime(): boolean {
  return typeof window !== "undefined" && Boolean((window as Window & { go?: unknown }).go);
}

function requireNativeRuntime(): void {
  if (!isNativeRuntime()) {
    throw new Error("The Agentic App Studio desktop application is required.");
  }
}

export async function listProjects(): Promise<Project[]> {
  requireNativeRuntime();
  return ListProjects();
}

export async function createProject(name: string, brief: string): Promise<Project> {
  requireNativeRuntime();
  return CreateProject(name, brief);
}

export async function health(): Promise<StudioHealth> {
  requireNativeRuntime();
  return Health();
}

export async function createRevision(projectId: string, snapshot: string): Promise<Revision> {
  requireNativeRuntime();
  return CreateRevision(projectId, snapshot);
}

export async function listRevisions(projectId: string): Promise<Revision[]> {
  requireNativeRuntime();
  return ListRevisions(projectId);
}

export async function compileRevision(projectId: string, revisionId: string): Promise<CompiledPlan> {
  requireNativeRuntime();
  return CompileRevision(projectId, revisionId);
}

export async function listAdapters(): Promise<AdapterManifest[]> {
  requireNativeRuntime();
  return AvailableAdapters();
}

export async function evaluateSchedule(projectId: string, revisionId: string, adapter: AdapterManifest): Promise<ScheduleResult> {
  const capabilities = adapter.capabilities;
  requireNativeRuntime();
  return EvaluateSchedule(projectId, revisionId, adapter.id, capabilities);
}

export async function executeCard(card: CardSpec, repository: string, worktreeRoot: string, baseRef: string, adapterId: string): Promise<RunResult> {
  requireNativeRuntime();
  return ExecuteCard(JSON.stringify(card), repository, worktreeRoot, baseRef, adapterId);
}

export async function startTerminal(command: string, directory: string): Promise<TerminalInfo> {
  requireNativeRuntime();
  if (!command.trim()) throw new Error("Choose a trusted agent command before opening a terminal");
  return StartTerminal(command, [], directory);
}

export async function listTerminalEvents(id: string): Promise<TerminalEvent[]> {
  requireNativeRuntime();
  return TerminalEvents(id);
}

export async function writeTerminalInput(id: string, input: string): Promise<void> {
  requireNativeRuntime();
  return WriteTerminalInput(id, input);
}

export async function interruptTerminal(id: string): Promise<void> {
  requireNativeRuntime();
  return InterruptTerminal(id);
}

export async function closeTerminal(id: string): Promise<void> {
  requireNativeRuntime();
  return CloseTerminal(id);
}

export async function defaultWorkspaceDirectory(): Promise<string> {
  requireNativeRuntime();
  return DefaultWorkspaceDirectory();
}

export async function selectWorkspaceDirectory(): Promise<string> {
  requireNativeRuntime();
  return SelectWorkspaceDirectory();
}
