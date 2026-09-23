import { FormEvent, useEffect, useRef, useState } from "react";
import { listTerminalEvents, selectWorkspaceDirectory, startTerminal, writeTerminalInput, interruptTerminal, closeTerminal, TerminalEvent, TerminalInfo } from "../studio-api";

export default function TerminalWorkspace() {
  const [command, setCommand] = useState("");
  const [directory, setDirectory] = useState("");
  const [input, setInput] = useState("");
  const [terminal, setTerminal] = useState<TerminalInfo | null>(null);
  const [events, setEvents] = useState<TerminalEvent[]>([]);
  const [feedback, setFeedback] = useState<string | null>(null);
  const eventCount = useRef(0);

  useEffect(() => {
    if (!terminal) return;
    const interval = window.setInterval(() => {
      void listTerminalEvents(terminal.id)
        .then((nextEvents) => {
          if (nextEvents.length === eventCount.current) return;
          eventCount.current = nextEvents.length;
          setEvents(nextEvents);
        })
        .catch(() => setFeedback("Unable to read terminal events."));
    }, 250);
    return () => window.clearInterval(interval);
  }, [terminal]);

  async function handleStart(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!command.trim()) {
      setFeedback("Command is required.");
      return;
    }
    setFeedback(null);
    setEvents([]);
    eventCount.current = 0;
    try {
      const selectedDirectory = directory.trim() || await selectWorkspaceDirectory();
      if (!selectedDirectory) {
        setFeedback("Choose a workspace folder to continue.");
        return;
      }
      setDirectory(selectedDirectory);
      setTerminal(await startTerminal(command.trim(), selectedDirectory));
    } catch (error) {
      setFeedback(error instanceof Error ? error.message : "Unable to start terminal.");
    }
  }

  async function handleInput(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!terminal || !input) return;
    try {
      await writeTerminalInput(terminal.id, input);
      setInput("");
    } catch {
      setFeedback("Unable to send input to terminal.");
    }
  }

  async function handleInterrupt() {
    if (!terminal) return;
    await interruptTerminal(terminal.id).catch(() => setFeedback("Unable to interrupt terminal."));
  }

  async function handleClose() {
    if (!terminal) return;
    await closeTerminal(terminal.id).catch(() => setFeedback("Unable to close terminal."));
    setTerminal(null);
    setEvents([]);
    eventCount.current = 0;
  }

  return (
    <section className="terminal-workspace" aria-labelledby="terminal-workspace-title">
      <div className="terminal-header">
        <div>
          <p className="wizard-kicker">Execution workspace</p>
          <h2 id="terminal-workspace-title">Live agent terminal</h2>
        </div>
        {terminal && <span className="terminal-state">{terminal.state}</span>}
      </div>

      {!terminal ? (
        <form className="terminal-start-form" onSubmit={handleStart}>
          <label>
            Agent command
            <input value={command} onChange={(event) => setCommand(event.target.value)} placeholder="Path or command for a trusted agent" autoFocus />
          </label>
          <div className="workspace-picker">
            <span className="workspace-label">Working directory</span>
            <p className="workspace-path">{directory || "No folder selected yet"}</p>
            <button type="button" className="secondary-action" onClick={() => void selectWorkspaceDirectory().then(setDirectory).catch((error) => setFeedback(error instanceof Error ? error.message : "Unable to choose a folder."))}>
              Choose folder
            </button>
          </div>
          <button type="submit" className="start-button">Open terminal</button>
        </form>
      ) : (
        <>
          <pre className="terminal-output" aria-live="polite">{events.map((event) => event.data).join("") || "Waiting for terminal output…"}</pre>
          <form className="terminal-input-form" onSubmit={handleInput}>
            <label htmlFor="terminal-input">Send input</label>
            <textarea id="terminal-input" value={input} onChange={(event) => setInput(event.target.value)} rows={3} placeholder="Type input for the agent…" />
            <div className="terminal-actions">
              <button type="submit" className="start-button">Send</button>
              <button type="button" className="secondary-action" onClick={() => void handleInterrupt()}>Ctrl+C</button>
              <button type="button" className="secondary-action" onClick={() => void handleClose()}>Close</button>
            </div>
          </form>
        </>
      )}
      {feedback && <p className="upload-error" role="alert">{feedback}</p>}
    </section>
  );
}
