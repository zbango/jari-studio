import { FormEvent, useEffect, useState } from "react";
import {
  createProject,
  compileRevision,
  createRevision,
  evaluateSchedule,
  executeCard,
  health,
  listAdapters,
  isNativeRuntime,
  listProjects,
  listRevisions,
  Project,
  CompiledPlan,
  ScheduleResult,
  AdapterManifest,
  RunResult,
  Revision,
  StudioHealth,
} from "./studio-api";
import MyStudioPage from "./components/MyStudioPage";
import TerminalWorkspace from "./components/TerminalWorkspace";
import TerminalTestPage from "./components/TerminalTestPage";

type ViewId = "my-studio" | "projects" | "brief" | "domain" | "plan" | "terminal" | "terminal-test" | "evidence";

type ViewDefinition = {
  id: ViewId;
  label: string;
  eyebrow: string;
  title: string;
  summary: string;
  bullets: string[];
};

type SnapshotDraft = {
  modules?: Array<{ name?: string }>;
  entities?: Array<{ name?: string; fields?: Array<{ name?: string }> }>;
  relations?: Array<{ fromEntity?: string; toEntity?: string; kind?: string }>;
  validations?: Array<{ targetId?: string; rule?: string }>;
  gaps?: Array<{ description?: string; severity?: string }>;
};

const views: ViewDefinition[] = [
  {
    id: "my-studio",
    label: "My Studio",
    eyebrow: "Workspace",
    title: "My Studio",
    summary: "Let's build something cool!",
    bullets: [],
  },
  {
    id: "projects",
    label: "Projects",
    eyebrow: "Studio Entry",
    title: "Local-first project control",
    summary:
      "Create and reopen Studio projects backed by local SQLite, versioned assets and Product Graph revisions.",
    bullets: [
      "Project metadata and local storage health",
      "Recent revisions and runtime diagnostics",
      "Workspace authority and adapter readiness",
    ],
  },
  {
    id: "brief",
    label: "Brief",
    eyebrow: "Discovery",
    title: "Business intent before execution",
    summary:
      "Capture the brief, optional screenshots and unresolved gaps before Cards are compiled or any agent is allowed to act.",
    bullets: [
      "Business goals, scope limits and operator notes",
      "Optional visual evidence and provenance",
      "Decision states: observed, inferred, proposed, confirmed",
    ],
  },
  {
    id: "plan",
    label: "Plan",
    eyebrow: "Orchestration",
    title: "Cards, sequencing and gates",
    summary:
      "Translate approved requirements into Cards with dependency locks, capability requirements and verification gates.",
    bullets: [
      "Card readiness, risk budget and semantic conflicts",
      "Adapter capability matching and execution policy",
      "Merge-verified workflow ending at approved evidence",
    ],
  },
  {
    id: "domain",
    label: "Domain",
    eyebrow: "Product Graph",
    title: "Entities before implementation",
    summary:
      "Capture the business objects and fields that downstream Cards will use without letting visual guesses become schema silently.",
    bullets: [
      "Entities and fields with explicit provenance",
      "Relations and validations as the next graph layer",
      "Proposed revisions remain reviewable before execution",
    ],
  },
  {
    id: "evidence",
    label: "Evidence",
    eyebrow: "Verification",
    title: "Traceability without guesswork",
    summary:
      "Link acceptance criteria, diffs, runs, artifacts and ledger entries so every approved change can be audited.",
    bullets: [
      "Verifier outputs and artifact hashes",
      "Human approval checkpoints",
      "Ledger navigation from intent to merge",
    ],
  },
  {
    id: "terminal",
    label: "Terminal",
    eyebrow: "Execution",
    title: "Live agent terminal",
    summary: "Open and control a real terminal session for a coding agent.",
    bullets: [
      "Send input and observe output in real time",
      "Interrupt or close the running process",
      "Keep the session tied to its working directory",
    ],
  },
  {
    id: "terminal-test",
    label: "Terminal Test",
    eyebrow: "Diagnostics",
    title: "Terminal test workspace",
    summary: "Test real agent terminal input and output independently from My Studio.",
    bullets: [
      "Launch a local command in its working directory",
      "Inspect output and send input interactively",
      "Interrupt or close the session safely",
    ],
  },
];

const metrics = [
  { label: "Primary Authority", value: "Desktop Studio" },
  { label: "Current Phase", value: "Foundation / UX bootstrap" },
  { label: "Execution Scope", value: "Merge verified" },
];

export default function App() {
  const [activeView, setActiveView] = useState<ViewId>("projects");
  const [projects, setProjects] = useState<Project[]>([]);
  const [selectedProjectId, setSelectedProjectId] = useState<string | null>(null);
  const [revisions, setRevisions] = useState<Revision[]>([]);
  const [compiledPlan, setCompiledPlan] = useState<CompiledPlan | null>(null);
  const [schedule, setSchedule] = useState<ScheduleResult | null>(null);
  const [adapters, setAdapters] = useState<AdapterManifest[]>([]);
  const [selectedAdapterId, setSelectedAdapterId] = useState("codex");
  const [repositoryPath, setRepositoryPath] = useState("");
  const [worktreeRoot, setWorktreeRoot] = useState("");
  const [baseRef, setBaseRef] = useState("HEAD");
  const [runningCardId, setRunningCardId] = useState<string | null>(null);
  const [runResult, setRunResult] = useState<RunResult | null>(null);
  const [briefDraft, setBriefDraft] = useState("");
  const [moduleDrafts, setModuleDrafts] = useState<string[]>([]);
  const [moduleInput, setModuleInput] = useState("");
  const [entityDrafts, setEntityDrafts] = useState<Array<{ name: string; fields: string[] }>>([]);
  const [entityInput, setEntityInput] = useState("");
  const [fieldInputs, setFieldInputs] = useState<Record<string, string>>({});
  const [relationDrafts, setRelationDrafts] = useState<Array<{ fromEntity: string; toEntity: string; kind: string }>>([]);
  const [relationInput, setRelationInput] = useState({ fromEntity: "", toEntity: "", kind: "one-to-many" });
  const [validationDrafts, setValidationDrafts] = useState<Array<{ targetId: string; rule: string }>>([]);
  const [validationInput, setValidationInput] = useState({ targetId: "", rule: "" });
  const [gapDrafts, setGapDrafts] = useState<Array<{ description: string; severity: string }>>([]);
  const [gapInput, setGapInput] = useState({ description: "", severity: "blocking" });
  const [studioHealth, setStudioHealth] = useState<StudioHealth | null>(null);
  const [projectName, setProjectName] = useState("");
  const [projectBrief, setProjectBrief] = useState("");
  const [feedback, setFeedback] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isSaving, setIsSaving] = useState(false);
  const [isSavingRevision, setIsSavingRevision] = useState(false);
  const [isCompiling, setIsCompiling] = useState(false);
  const current = views.find((view) => view.id === activeView) ?? views[0];
  const selectedProject = projects.find((project) => project.id === selectedProjectId) ?? null;

  useEffect(() => {
    if (!isNativeRuntime()) {
      setStudioHealth({ status: "desktop_required", runtime: "browser-preview" });
      setIsLoading(false);
      return;
    }
    void Promise.all([listProjects(), health(), listAdapters()])
      .then(([loadedProjects, loadedHealth, loadedAdapters]) => {
        setProjects(loadedProjects);
        setSelectedProjectId(loadedProjects[0]?.id ?? null);
        setStudioHealth(loadedHealth);
        setAdapters(loadedAdapters);
        if (loadedAdapters.length > 0 && !loadedAdapters.some((adapter) => adapter.id === selectedAdapterId)) {
          setSelectedAdapterId(loadedAdapters[0].id);
        }
      })
      .catch(() => setFeedback("No se pudo cargar el estado local del Studio."))
      .finally(() => setIsLoading(false));
  }, []);

  useEffect(() => {
    if (!selectedProject) {
      setBriefDraft("");
    setRevisions([]);
      setCompiledPlan(null);
      setSchedule(null);
      return;
    }
    setBriefDraft(selectedProject.brief ?? "");
    setModuleDrafts([]);
    setModuleInput("");
    setEntityDrafts([]);
    setEntityInput("");
    setFieldInputs({});
    setRelationDrafts([]);
    setRelationInput({ fromEntity: "", toEntity: "", kind: "one-to-many" });
    setValidationDrafts([]);
    setValidationInput({ targetId: "", rule: "" });
    setGapDrafts([]);
    setGapInput({ description: "", severity: "blocking" });
    setCompiledPlan(null);
    setSchedule(null);
    void listRevisions(selectedProject.id)
      .then((loadedRevisions) => {
        setRevisions(loadedRevisions);
        const latestSnapshot = loadedRevisions.at(-1)?.snapshot;
        if (!latestSnapshot) {
          return;
        }
        try {
          const snapshot = JSON.parse(latestSnapshot) as SnapshotDraft;
          setModuleDrafts((snapshot.modules ?? []).flatMap((module) => module.name ? [module.name] : []));
          setEntityDrafts((snapshot.entities ?? []).flatMap((entity) => entity.name ? [{
            name: entity.name,
            fields: (entity.fields ?? []).flatMap((field) => field.name ? [field.name] : []),
          }] : []));
          setRelationDrafts((snapshot.relations ?? []).flatMap((relation) => relation.fromEntity && relation.toEntity ? [{
            fromEntity: relation.fromEntity,
            toEntity: relation.toEntity,
            kind: relation.kind ?? "one-to-many",
          }] : []));
          setValidationDrafts((snapshot.validations ?? []).flatMap((validation) => validation.targetId && validation.rule ? [{
            targetId: validation.targetId,
            rule: validation.rule,
          }] : []));
          setGapDrafts((snapshot.gaps ?? []).flatMap((gap) => gap.description ? [{
            description: gap.description,
            severity: gap.severity ?? "blocking",
          }] : []));
        } catch {
          setFeedback("La última revisión no contiene un snapshot legible.");
        }
      })
      .catch(() => setRevisions([]));
  }, [selectedProject]);

  async function handleCreateProject(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const name = projectName.trim();
    const brief = projectBrief.trim();
    if (!name || !brief) {
      setFeedback("Nombre y brief son obligatorios.");
      return;
    }

    setIsSaving(true);
    setFeedback(null);
    try {
      const project = await createProject(name, brief);
      setProjects((currentProjects) => [...currentProjects, project]);
      setProjectName("");
      setProjectBrief("");
      setFeedback(`Proyecto “${project.name}” creado en la autoridad local.`);
    } catch {
      setFeedback("No se pudo crear el proyecto.");
    } finally {
      setIsSaving(false);
    }
  }

  async function handleSaveBrief(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!selectedProject || !briefDraft.trim()) {
      setFeedback("Selecciona un proyecto y escribe un brief.");
      return;
    }

    setIsSavingRevision(true);
    setFeedback(null);
    const snapshot = JSON.stringify({
      project: {
        id: selectedProject.id,
        name: selectedProject.name,
        brief: briefDraft.trim(),
      },
      modules: moduleDrafts.map((name, index) => ({
        id: `${selectedProject.id}.module.${index + 1}`,
        projectId: selectedProject.id,
        name,
        metadata: {
          status: "proposed",
          revision: 1,
          sources: [{ kind: "client_decision", ref: "studio.brief" }],
        },
      })),
      decisions: [],
      source: "studio.brief",
    });
    try {
      const revision = await createRevision(selectedProject.id, snapshot);
      setRevisions((currentRevisions) => [...currentRevisions, revision]);
      setFeedback(`Revisión ${revision.number} guardada en SQLite.`);
    } catch {
      setFeedback("No se pudo guardar la revisión del brief.");
    } finally {
      setIsSavingRevision(false);
    }
  }

  function handleAddModule() {
    const moduleName = moduleInput.trim();
    if (!moduleName || moduleDrafts.includes(moduleName)) {
      return;
    }
    setModuleDrafts((currentModules) => [...currentModules, moduleName]);
    setModuleInput("");
  }

  function handleAddEntity() {
    const entityName = entityInput.trim();
    if (!entityName || entityDrafts.some((entity) => entity.name === entityName)) {
      return;
    }
    setEntityDrafts((currentEntities) => [...currentEntities, { name: entityName, fields: [] }]);
    setEntityInput("");
  }

  function handleAddField(entityName: string) {
    const fieldName = (fieldInputs[entityName] ?? "").trim();
    if (!fieldName) {
      return;
    }
    setEntityDrafts((currentEntities) => currentEntities.map((entity) => (
      entity.name === entityName && !entity.fields.includes(fieldName)
        ? { ...entity, fields: [...entity.fields, fieldName] }
        : entity
    )));
    setFieldInputs((currentInputs) => ({ ...currentInputs, [entityName]: "" }));
  }

  async function handleSaveDomain(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!selectedProject || entityDrafts.length === 0) {
      setFeedback("Selecciona un proyecto y añade al menos una entidad.");
      return;
    }
    setIsSavingRevision(true);
    setFeedback(null);
    const snapshot = JSON.stringify({
      project: { id: selectedProject.id, name: selectedProject.name, brief: briefDraft.trim() },
      modules: moduleDrafts.map((name, index) => ({
        id: `${selectedProject.id}.module.${index + 1}`,
        projectId: selectedProject.id,
        name,
        metadata: { status: "proposed", revision: 1, sources: [{ kind: "client_decision", ref: "studio.domain" }] },
      })),
      entities: entityDrafts.map((entity, entityIndex) => ({
        id: `${selectedProject.id}.entity.${entityIndex + 1}`,
        name: entity.name,
        fields: entity.fields.map((field, fieldIndex) => ({
          id: `${selectedProject.id}.entity.${entityIndex + 1}.field.${fieldIndex + 1}`,
          name: field,
          type: "string",
          metadata: { status: "proposed", revision: 1, sources: [{ kind: "client_decision", ref: "studio.domain" }] },
        })),
        metadata: { status: "proposed", revision: 1, sources: [{ kind: "client_decision", ref: "studio.domain" }] },
      })),
      relations: relationDrafts.map((relation, index) => ({
        id: `${selectedProject.id}.relation.${index + 1}`,
        fromEntity: relation.fromEntity,
        toEntity: relation.toEntity,
        kind: relation.kind,
        metadata: { status: "proposed", revision: 1, sources: [{ kind: "client_decision", ref: "studio.domain" }] },
      })),
      validations: validationDrafts.map((validation, index) => ({
        id: `${selectedProject.id}.validation.${index + 1}`,
        targetId: validation.targetId,
        rule: validation.rule,
        metadata: { status: "proposed", revision: 1, sources: [{ kind: "client_decision", ref: "studio.domain" }] },
      })),
      gaps: gapDrafts.map((gap, index) => ({
        id: `${selectedProject.id}.gap.${index + 1}`,
        type: "domain",
        severity: gap.severity,
        description: gap.description,
        status: "proposed",
      })),
      source: "studio.domain",
    });
    try {
      const revision = await createRevision(selectedProject.id, snapshot);
      setRevisions((currentRevisions) => [...currentRevisions, revision]);
      setFeedback(`Revisión de dominio ${revision.number} guardada en SQLite.`);
    } catch {
      setFeedback("No se pudo guardar la revisión de dominio.");
    } finally {
      setIsSavingRevision(false);
    }
  }

  async function handleCompilePlan() {
    if (!selectedProject || revisions.length === 0) {
      setFeedback("Guarda una revisión antes de compilar el plan.");
      return;
    }
    setIsCompiling(true);
    setFeedback(null);
    try {
      const latestRevision = revisions[revisions.length - 1];
      const plan = await compileRevision(selectedProject.id, latestRevision.id);
      setCompiledPlan(plan);
      const adapter = adapters.find((candidate) => candidate.id === selectedAdapterId) ?? adapters[0];
      if (!adapter) {
        throw new Error("no adapter available");
      }
      setSchedule(await evaluateSchedule(selectedProject.id, latestRevision.id, adapter));
      setFeedback(`Plan compilado desde la revisión ${latestRevision.number}.`);
    } catch {
      setFeedback("No se pudo compilar el plan. Revisa los gaps blocking.");
    } finally {
      setIsCompiling(false);
    }
  }

  async function handleRunCard(card: CompiledPlan["cards"][number]) {
    if (!repositoryPath.trim() || !worktreeRoot.trim() || !selectedAdapterId) {
      setFeedback("Indica repositorio, raíz de worktrees y adapter antes de ejecutar.");
      return;
    }
    setRunningCardId(card.id);
    setRunResult(null);
    setFeedback(null);
    try {
      const result = await executeCard(card, repositoryPath.trim(), worktreeRoot.trim(), baseRef.trim() || "HEAD", selectedAdapterId);
      setRunResult(result);
      setFeedback(`Card ${card.id} terminó en estado ${result.state}.`);
    } catch (error) {
      setFeedback(error instanceof Error ? error.message : "No se pudo ejecutar la Card.");
    } finally {
      setRunningCardId(null);
    }
  }

  function handleAddRelation() {
    const { fromEntity, toEntity, kind } = relationInput;
    if (!fromEntity || !toEntity || fromEntity === toEntity) return;
    setRelationDrafts((currentRelations) => [...currentRelations, { fromEntity, toEntity, kind }]);
    setRelationInput({ fromEntity: "", toEntity: "", kind: "one-to-many" });
  }

  function handleAddValidation() {
    if (!validationInput.targetId.trim() || !validationInput.rule.trim()) return;
    setValidationDrafts((currentValidations) => [...currentValidations, { targetId: validationInput.targetId.trim(), rule: validationInput.rule.trim() }]);
    setValidationInput({ targetId: "", rule: "" });
  }

  function handleAddGap() {
    if (!gapInput.description.trim()) return;
    setGapDrafts((currentGaps) => [...currentGaps, { description: gapInput.description.trim(), severity: gapInput.severity }]);
    setGapInput({ description: "", severity: "blocking" });
  }

  return (
    <div className="app-shell">
      <aside className="sidebar" aria-label="Studio sections">
        <div className="brand-block">
          <p className="kicker">Agentic App Studio</p>
          <h1>Studio shell</h1>
          <p className="brand-copy">
            Wails desktop frontend bootstrap for the local-first orchestration
            surface.
          </p>
        </div>

        <nav className="nav-list" aria-label="Primary navigation">
          {views.map((view) => (
            <button
              key={view.id}
              type="button"
              className={view.id === activeView ? "nav-item active" : "nav-item"}
              onClick={() => setActiveView(view.id)}
              aria-current={view.id === activeView ? "page" : undefined}
            >
              <span>{view.label}</span>
              <small>{view.eyebrow}</small>
            </button>
          ))}
        </nav>
      </aside>

      <main className={activeView === "my-studio" ? "content my-studio-page" : "content"} aria-live="polite">
        {studioHealth?.status === "desktop_required" && (
          <section className="desktop-required" role="status">
            <strong>Agentic App Studio runs as a desktop application.</strong>
            <span>Open the Wails host to access SQLite, native folders and terminal sessions.</span>
          </section>
        )}
        {activeView === "my-studio" && (
          <MyStudioPage />
        )}
        <section className="hero-card">
          <div>
            <p className="kicker">{current.eyebrow}</p>
            <h2>{current.title}</h2>
            <p className="summary">{current.summary}</p>
          </div>

          <div className="hero-status">
            <span className="status-pill">Offline-ready</span>
            <span className="status-pill">No cloud dependency</span>
            <span className="status-pill">Read model aware</span>
          </div>
        </section>

        <section className="panel-grid" aria-label="Studio overview">
          <article className="panel">
            <h3>What this screen must support</h3>
            <ul>
              {current.bullets.map((bullet) => (
                <li key={bullet}>{bullet}</li>
              ))}
            </ul>
          </article>

          <article className="panel">
            <h3>Foundation metrics</h3>
            <dl className="metrics">
              {[...metrics, { label: "Runtime", value: studioHealth?.runtime ?? "loading" }].map((metric) => (
                <div key={metric.label}>
                  <dt>{metric.label}</dt>
                  <dd>{metric.value}</dd>
                </div>
              ))}
            </dl>
          </article>
        </section>

        {activeView === "projects" && (
          <section className="panel project-panel" aria-labelledby="projects-heading">
            <div className="panel-heading">
              <div>
                <p className="kicker">Projects</p>
                <h3 id="projects-heading">Local project authority</h3>
              </div>
              <span className="count-badge">{isLoading ? "…" : projects.length}</span>
            </div>

            <form className="project-form" onSubmit={handleCreateProject}>
              <label>
                Nombre del proyecto
                <input value={projectName} onChange={(event) => setProjectName(event.target.value)} placeholder="Ej. Gestión de membresías" />
              </label>
              <label>
                Brief
                <textarea value={projectBrief} onChange={(event) => setProjectBrief(event.target.value)} placeholder="Describe el problema y el resultado esperado." rows={3} />
              </label>
              <button type="submit" className="primary-action" disabled={isSaving}>
                {isSaving ? "Creando…" : "Crear proyecto"}
              </button>
            </form>

            {feedback && <p className="feedback" role="status">{feedback}</p>}

            <div className="project-list" aria-label="Projects created">
              {projects.length === 0 && !isLoading ? (
                <p className="empty-state">Aún no hay proyectos. Crea el primero para iniciar su Product Graph.</p>
              ) : (
                projects.map((project) => (
                  <article
                    className={project.id === selectedProjectId ? "project-row selected" : "project-row"}
                    key={project.id}
                  >
                    <div>
                      <strong>{project.name}</strong>
                      <p>{project.brief}</p>
                    </div>
                    <button
                      type="button"
                      className="row-action"
                      onClick={() => {
                        setSelectedProjectId(project.id);
                        setActiveView("brief");
                      }}
                    >
                      Abrir
                    </button>
                  </article>
                ))
              )}
            </div>
          </section>
        )}

        {activeView === "brief" && (
          <section className="panel project-panel" aria-labelledby="brief-heading">
            <div className="panel-heading">
              <div>
                <p className="kicker">Brief</p>
                <h3 id="brief-heading">Approved product intent</h3>
              </div>
              <span className="count-badge">{revisions.length}</span>
            </div>

            {!selectedProject ? (
              <p className="empty-state">Crea o selecciona un proyecto antes de editar su brief.</p>
            ) : (
              <form className="project-form" onSubmit={handleSaveBrief}>
                <p className="selected-project">Proyecto: <strong>{selectedProject.name}</strong></p>
                <label>
                  Brief del producto
                  <textarea
                    value={briefDraft}
                    onChange={(event) => setBriefDraft(event.target.value)}
                    rows={7}
                    aria-describedby="brief-help"
                  />
                </label>
                <div className="module-editor">
                  <label htmlFor="module-input">Módulos iniciales</label>
                  <div className="inline-form">
                    <input
                      id="module-input"
                      value={moduleInput}
                      onChange={(event) => setModuleInput(event.target.value)}
                      onKeyDown={(event) => {
                        if (event.key === "Enter") {
                          event.preventDefault();
                          handleAddModule();
                        }
                      }}
                      placeholder="Ej. Clientes"
                    />
                    <button type="button" className="row-action" onClick={handleAddModule}>Agregar</button>
                  </div>
                  {moduleDrafts.length > 0 && (
                    <ul className="module-list" aria-label="Módulos propuestos">
                      {moduleDrafts.map((moduleName) => <li key={moduleName}>{moduleName}</li>)}
                    </ul>
                  )}
                </div>
                <p id="brief-help" className="field-help">
                  Guardar crea una nueva revisión propuesta del Product Graph. No ejecuta agentes.
                </p>
                <button type="submit" className="primary-action" disabled={isSavingRevision}>
                  {isSavingRevision ? "Guardando…" : "Guardar revisión"}
                </button>
              </form>
            )}

            {feedback && <p className="feedback" role="status">{feedback}</p>}

            <div className="revision-list" aria-label="Product Graph revisions">
              {revisions.length === 0 ? (
                <p className="empty-state">Aún no hay revisiones guardadas.</p>
              ) : (
                revisions.map((revision) => (
                  <article className="revision-row" key={revision.id}>
                    <strong>Revision {revision.number}</strong>
                    <span>{revision.status}</span>
                  </article>
                ))
              )}
            </div>
          </section>
        )}

        {activeView === "domain" && (
          <section className="panel project-panel" aria-labelledby="domain-heading">
            <div className="panel-heading">
              <div>
                <p className="kicker">Domain</p>
                <h3 id="domain-heading">Proposed business entities</h3>
              </div>
              <span className="count-badge">{entityDrafts.length}</span>
            </div>

            {!selectedProject ? (
              <p className="empty-state">Crea o selecciona un proyecto antes de modelar su dominio.</p>
            ) : (
              <form className="project-form" onSubmit={handleSaveDomain}>
                <p className="selected-project">Proyecto: <strong>{selectedProject.name}</strong></p>
                <div className="inline-form">
                  <input
                    value={entityInput}
                    onChange={(event) => setEntityInput(event.target.value)}
                    onKeyDown={(event) => {
                      if (event.key === "Enter") {
                        event.preventDefault();
                        handleAddEntity();
                      }
                    }}
                    placeholder="Ej. Member"
                    aria-label="Nueva entidad"
                  />
                  <button type="button" className="row-action" onClick={handleAddEntity}>Agregar entidad</button>
                </div>
                <div className="entity-list" aria-label="Entidades propuestas">
                  {entityDrafts.map((entity) => (
                    <article className="entity-card" key={entity.name}>
                      <strong>{entity.name}</strong>
                      <div className="inline-form">
                        <input
                          value={fieldInputs[entity.name] ?? ""}
                          onChange={(event) => setFieldInputs((currentInputs) => ({ ...currentInputs, [entity.name]: event.target.value }))}
                          onKeyDown={(event) => {
                            if (event.key === "Enter") {
                              event.preventDefault();
                              handleAddField(entity.name);
                            }
                          }}
                          placeholder="Campo, ej. email"
                          aria-label={`Nuevo campo para ${entity.name}`}
                        />
                        <button type="button" className="row-action" onClick={() => handleAddField(entity.name)}>Agregar campo</button>
                      </div>
                      {entity.fields.length > 0 && <ul className="module-list">{entity.fields.map((field) => <li key={field}>{field}</li>)}</ul>}
                    </article>
                  ))}
                </div>
                <div className="domain-subsection">
                  <label>Relaciones</label>
                  <div className="inline-form">
                    <select value={relationInput.fromEntity} onChange={(event) => setRelationInput((currentInput) => ({ ...currentInput, fromEntity: event.target.value }))} aria-label="Entidad origen">
                      <option value="">Origen</option>
                      {entityDrafts.map((entity) => <option key={entity.name} value={entity.name}>{entity.name}</option>)}
                    </select>
                    <select value={relationInput.toEntity} onChange={(event) => setRelationInput((currentInput) => ({ ...currentInput, toEntity: event.target.value }))} aria-label="Entidad destino">
                      <option value="">Destino</option>
                      {entityDrafts.map((entity) => <option key={entity.name} value={entity.name}>{entity.name}</option>)}
                    </select>
                    <select value={relationInput.kind} onChange={(event) => setRelationInput((currentInput) => ({ ...currentInput, kind: event.target.value }))} aria-label="Tipo de relación">
                      <option value="one-to-one">1:1</option>
                      <option value="one-to-many">1:N</option>
                      <option value="many-to-many">N:M</option>
                    </select>
                    <button type="button" className="row-action" onClick={handleAddRelation}>Agregar</button>
                  </div>
                  {relationDrafts.length > 0 && <ul className="domain-list">{relationDrafts.map((relation, index) => <li key={`${relation.fromEntity}-${relation.toEntity}-${index}`}>{relation.fromEntity} → {relation.toEntity} ({relation.kind})</li>)}</ul>}
                </div>
                <div className="domain-subsection">
                  <label htmlFor="validation-target">Validaciones</label>
                  <div className="inline-form">
                    <input id="validation-target" value={validationInput.targetId} onChange={(event) => setValidationInput((currentInput) => ({ ...currentInput, targetId: event.target.value }))} placeholder="Campo objetivo" />
                    <input value={validationInput.rule} onChange={(event) => setValidationInput((currentInput) => ({ ...currentInput, rule: event.target.value }))} placeholder="Regla, ej. required" aria-label="Regla de validación" />
                    <button type="button" className="row-action" onClick={handleAddValidation}>Agregar</button>
                  </div>
                  {validationDrafts.length > 0 && <ul className="domain-list">{validationDrafts.map((validation, index) => <li key={`${validation.targetId}-${index}`}>{validation.targetId}: {validation.rule}</li>)}</ul>}
                </div>
                <div className="domain-subsection">
                  <label htmlFor="gap-description">Gaps detectados</label>
                  <div className="inline-form">
                    <input id="gap-description" value={gapInput.description} onChange={(event) => setGapInput((currentInput) => ({ ...currentInput, description: event.target.value }))} placeholder="Decisión pendiente" />
                    <select value={gapInput.severity} onChange={(event) => setGapInput((currentInput) => ({ ...currentInput, severity: event.target.value }))} aria-label="Severidad del gap">
                      <option value="blocking">Blocking</option>
                      <option value="non_blocking">Non-blocking</option>
                    </select>
                    <button type="button" className="row-action" onClick={handleAddGap}>Agregar</button>
                  </div>
                  {gapDrafts.length > 0 && <ul className="domain-list">{gapDrafts.map((gap, index) => <li key={`${gap.description}-${index}`}>{gap.severity}: {gap.description}</li>)}</ul>}
                </div>
                <p className="field-help">Las entidades y campos se guardan como propuestas; todavía no generan migraciones ni ejecutan agentes.</p>
                <button type="submit" className="primary-action" disabled={isSavingRevision}>
                  {isSavingRevision ? "Guardando…" : "Guardar revisión de dominio"}
                </button>
              </form>
            )}

            {feedback && <p className="feedback" role="status">{feedback}</p>}
          </section>
        )}

        {activeView === "plan" && (
          <section className="panel project-panel" aria-labelledby="plan-heading">
            <div className="panel-heading">
              <div>
                <p className="kicker">Plan</p>
                <h3 id="plan-heading">Compiled execution cards</h3>
              </div>
              <button type="button" className="primary-action" onClick={() => void handleCompilePlan()} disabled={isCompiling}>
                {isCompiling ? "Compilando…" : "Compilar plan"}
              </button>
            </div>
            {!selectedProject ? (
              <p className="empty-state">Crea o selecciona un proyecto antes de compilar Cards.</p>
            ) : !compiledPlan ? (
              <p className="empty-state">Selecciona “Compilar plan” para transformar la última revisión en Requirements, AC y Cards.</p>
            ) : (
              <div className="compiled-plan" aria-label="Compiled cards">
                <label className="adapter-picker">
                  Adapter de ejecución
                  <select value={selectedAdapterId} onChange={(event) => setSelectedAdapterId(event.target.value)}>
                    {adapters.map((adapter) => <option value={adapter.id} key={adapter.id}>{adapter.id} · {adapter.transport}</option>)}
                  </select>
                </label>
                <div className="plan-summary">
                  <span>{compiledPlan.requirements.length} requirements</span>
                  <span>{compiledPlan.acceptanceCriteria.length} acceptance criteria</span>
                  <span>{compiledPlan.cards.length} cards</span>
                </div>
                <div className="execution-config">
                  <label>Repositorio local<input value={repositoryPath} onChange={(event) => setRepositoryPath(event.target.value)} placeholder="/ruta/al/repositorio" /></label>
                  <label>Raíz de worktrees<input value={worktreeRoot} onChange={(event) => setWorktreeRoot(event.target.value)} placeholder="/ruta/worktrees" /></label>
                  <label>Base ref<input value={baseRef} onChange={(event) => setBaseRef(event.target.value)} placeholder="HEAD o main" /></label>
                </div>
                {compiledPlan.blockingGaps && compiledPlan.blockingGaps.length > 0 && (
                  <p className="feedback warning" role="alert">Hay {compiledPlan.blockingGaps.length} gaps blocking; el plan no puede ejecutarse.</p>
                )}
                {compiledPlan.cards.map((card) => (
                  <article className="card-row" key={card.id}>
                    <div>
                      <strong>{card.goal}</strong>
                      <p>{card.id}</p>
                      {schedule?.decisions.find((decision) => decision.cardId === card.id)?.reasons?.map((reason) => (
                        <small className="card-reason" key={reason}>{reason}</small>
                      ))}
                    </div>
                    <div className="card-actions">
                      <span className={`card-state ${schedule?.decisions.find((decision) => decision.cardId === card.id)?.state ?? card.state}`}>{schedule?.decisions.find((decision) => decision.cardId === card.id)?.state ?? card.state}</span>
                      {(schedule?.decisions.find((decision) => decision.cardId === card.id)?.state ?? card.state) === "ready" && (
                        <button type="button" className="row-action" onClick={() => void handleRunCard(card)} disabled={runningCardId !== null}>
                          {runningCardId === card.id ? "Ejecutando…" : "Run Card"}
                        </button>
                      )}
                    </div>
                  </article>
                ))}
                {runResult && <pre className="run-result" aria-label="Card execution result">{JSON.stringify(runResult, null, 2)}</pre>}
              </div>
            )}
            {feedback && <p className="feedback" role="status">{feedback}</p>}
          </section>
        )}

        {activeView === "terminal" && <TerminalWorkspace />}
        {activeView === "terminal-test" && <TerminalTestPage />}

        <section className="panel bridge-panel">
          <div>
            <h3>Next frontend integrations</h3>
            <p>
              Projects, Product Graph revisions and terminal sessions are
              available through the native Go/Wails host.
            </p>
          </div>

          <span className="status-pill">{studioHealth?.status ?? "Detecting runtime…"}</span>
        </section>
      </main>
    </div>
  );
}
