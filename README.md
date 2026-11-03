# Site Health Monitoring

Local desktop application for monitoring system health information.

## Current Status

The project currently contains a basic Electron desktop application that opens a local monitoring window.

## Requirements

Before running the application, make sure the following are installed:

- Node.js
- npm

## Run the Local Application

Open a terminal in the project directory and move into the local application folder:

```bash
cd local-app
```

Install the project dependencies:

```bash
npm install
```

Start the Electron application:

```bash
npm run dev
```

A desktop window should open and display the Site Health Monitoring interface.

## Project Structure

```text
site-health-monitoring/
├── local-app/
│   ├── backend/
│   ├── frontend/
│   │   └── index.html
│   ├── index.js
│   ├── package.json
│   └── package-lock.json
├── .gitignore
└── README.md
```

## Development

The project is currently under development. Future functionality will include local system data collection, data storage, alert evaluation, monitoring visualization, and cloud monitoring support.