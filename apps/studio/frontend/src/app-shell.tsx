export type StudioModule = {
  id: string;
  title: string;
  description: string;
  status: "planned" | "bootstrapped";
};

export const STUDIO_MODULES: StudioModule[] = [
  {
    id: "projects",
    title: "Projects",
    description: "Create and browse local-first Studio projects.",
    status: "bootstrapped"
  },
  {
    id: "brief",
    title: "Brief",
    description: "Capture intent, constraints and explicit business decisions.",
    status: "bootstrapped"
  },
  {
    id: "screenshots",
    title: "Screenshots",
    description: "Optional visual evidence and asset intake.",
    status: "planned"
  },
  {
    id: "domain",
    title: "Domain",
    description: "Review proposed modules, entities and rules.",
    status: "planned"
  },
  {
    id: "plan",
    title: "Plan",
    description: "Inspect Cards, dependencies and readiness.",
    status: "planned"
  },
  {
    id: "evidence",
    title: "Evidence",
    description: "Track verification artifacts and ledger links.",
    status: "planned"
  }
];

export function renderStudioShell(): string {
  return STUDIO_MODULES.map((module) => `${module.title}: ${module.status}`).join("\n");
}
