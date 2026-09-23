import { ChangeEvent, FormEvent, useEffect, useState } from "react";

type WizardStep = 0 | 1;

type MyStudioPageProps = {
  onSubmit?: (details: { oneLiner: string; description: string; image: File | null }) => void;
};

export default function MyStudioPage({ onSubmit }: MyStudioPageProps) {
  const [hasStarted, setHasStarted] = useState(false);
  const [step, setStep] = useState<WizardStep>(0);
  const [productOneLiner, setProductOneLiner] = useState("");
  const [productDescription, setProductDescription] = useState("");
  const [selectedImage, setSelectedImage] = useState<File | null>(null);
  const [previewUrl, setPreviewUrl] = useState<string | null>(null);
  const [uploadError, setUploadError] = useState<string | null>(null);

  useEffect(() => {
    if (!selectedImage) {
      setPreviewUrl(null);
      return;
    }

    let cancelled = false;
    let objectUrl: string | null = null;
    const image = selectedImage;

    async function preparePreview() {
      try {
        objectUrl = URL.createObjectURL(image);
        if (!cancelled) {
          setPreviewUrl(objectUrl);
          setUploadError(null);
        }
      } catch {
        setPreviewUrl(null);
        setUploadError("We couldn't preview this image. Please try another supported format.");
      }
    }

    void preparePreview();
    return () => {
      cancelled = true;
      if (objectUrl) URL.revokeObjectURL(objectUrl);
    };
  }, [selectedImage]);

  function handleImageChange(event: ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0] ?? null;
    if (!file) return;
    const isImage = file.type.startsWith("image/") && !/\.(heic|heif)$/i.test(file.name) && !/image\/hei[cf]/i.test(file.type);
    if (!isImage) {
      setSelectedImage(null);
      setPreviewUrl(null);
      setUploadError("Please choose an image file.");
      return;
    }
    setUploadError(null);
    setSelectedImage(file);
  }

  function handleContinue(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!productOneLiner.trim() || !productDescription.trim()) return;
    setStep(1);
  }

  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    onSubmit?.({ oneLiner: productOneLiner.trim(), description: productDescription.trim(), image: selectedImage });
  }

  return (
    <section className="empty-studio-page" aria-label="My Studio wizard">
      <div className="start-content wizard-content">
        {!hasStarted ? (
          <div className="wizard-step start-step" key="start-step">
            <p>Let's build something cool!</p>
            <button type="button" className="start-button" onClick={() => setHasStarted(true)}>Start</button>
          </div>
        ) : step === 0 ? (
          <div className="wizard-step" key="brief-step">
            <p className="wizard-kicker">Start a new project</p>
            <h2>Let's build something cool!</h2>
            <form className="wizard-form" onSubmit={handleContinue}>
              <label>
                What would you like to build in one sentence?
                <input
                  value={productOneLiner}
                  onChange={(event) => setProductOneLiner(event.target.value)}
                  placeholder="A simple way to manage my business"
                  autoFocus
                />
              </label>
              <label>
                Tell us briefly what it should do.
                <textarea
                  value={productDescription}
                  onChange={(event) => setProductDescription(event.target.value)}
                  placeholder="Describe the main problem it solves and who will use it."
                  rows={4}
                />
              </label>
              <button type="submit" className="start-button">Continue</button>
            </form>
          </div>
        ) : (
          <div className="wizard-step" key="upload-step">
            <p className="wizard-kicker">Step 2 of 2</p>
            <h2>Add a visual reference</h2>
            <p className="wizard-description">Upload a screen or image that captures the look and feel you have in mind.</p>
            <form className="upload-form" onSubmit={handleSubmit}>
              <label className="upload-dropzone" htmlFor="studio-image-upload">
                <span>Choose an image to preview</span>
                <small>PNG, JPG, GIF, SVG or WebP</small>
                <input id="studio-image-upload" type="file" accept="image/png,image/jpeg,image/gif,image/svg+xml,image/webp" onChange={handleImageChange} />
              </label>
              {uploadError && <p className="upload-error" role="alert">{uploadError}</p>}
              {previewUrl && selectedImage && (
                <figure className="image-preview">
                  <img src={previewUrl} alt={`Preview of ${selectedImage.name}`} />
                  <figcaption>{selectedImage.name}</figcaption>
                </figure>
              )}
              <div className="wizard-actions">
                <button type="button" className="secondary-action" onClick={() => setStep(0)}>Back</button>
                <button type="submit" className="start-button">Submit</button>
              </div>
            </form>
          </div>
        )}
      </div>
    </section>
  );
}
