# React frontend (source)

The React UI is intentionally supplied **without JSX or a Node.js build step**, so it can run offline inside the bundled Windows executable. This is real React code using `React.createElement`, class components and the React 16 UMD distribution. The source is `frontend/src/app.js` and the embedded copy is `backend/web/app.js`.

The UI makes HTTP requests to `/api/tasks`, `/api/stats`, and `/api/export.csv`. It supports the classic 2025 routes `/`, `/tasks/new`, `/tasks/:id/edit`, plus `/board` and `/stats`. Updating source: copy `frontend/src/app.js` into `backend/web/app.js` and rebuild Go to update the embedded executable.

The team's original frontend used Create React App, Axios, Bootstrap, React Router. This **independently written reconstruction** replaces those runtime dependencies with React UMD, native fetch, and custom Bootstrap-inspired CSS for a one-click offline startup. It preserves the main task workflow and the server's REST interface. Do not mistake it for the original friend's source.
