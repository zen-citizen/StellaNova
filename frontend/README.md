# StellaNova Frontend

This is a React application built with Vite and TypeScript.

## Tech Stack

- **React 19** - UI library
- **TypeScript** - Type safety
- **Vite** - Build tool and dev server
- **Tailwind CSS** - Styling with custom design tokens
- **ESLint** - Code linting
- **Prettier** - Code formatting
- **@vis.gl/react-google-maps** - Google Maps integration
- **MapLibre GL** - Map rendering

## Getting Started

### 1. Install Node Version

Use the Node version listed in the `.nvmrc` file or if you have nvm installed, run:

```bash
nvm use
```

### 2. Install Dependencies

```bash
npm install
```

### 3. Environment Variables

Create a `.env` file in the root directory with the following variables:

```env
VITE_GMAPS_API_KEY=your_google_maps_api_key_here
VITE_API_BASE_URL=https://api.stellanova.zencitizen.in
```

### 4. Run Development Server

```bash
npm run dev
```

The application will be available at `http://localhost:5173` (or the port Vite assigns).

## Available Scripts

- `npm run dev` - Start the development server
- `npm run build` - Build the application for production
- `npm run preview` - Preview the production build locally
- `npm run lint` - Run ESLint to check for code issues
- `npm run format` - Format code using Prettier

## Architecture

### Project Structure

```
src/
├── api.ts                 # API functions for fetching entities and addresses
├── App.tsx                # Main application component
├── main.tsx               # Application entry point with providers
├── styles.css             # Global styles and Tailwind design tokens
├── types.d.ts             # TypeScript type definitions
├── components/            # React components
│   ├── Accordion.tsx      # Accordion UI component
│   ├── Details.tsx        # Details view component
│   ├── Introduction.tsx   # Introduction/welcome view
│   ├── Map.tsx            # Map component with location markers
│   ├── Search.tsx         # Search input component
│   ├── Sidebar.tsx        # Main sidebar container
│   ├── SidebarContent.tsx # Sidebar content wrapper
│   ├── SidebarFooter.tsx  # Sidebar footer component
│   ├── SidebarHeader.tsx  # Sidebar header component
│   └── Table.tsx          # Table component for displaying entities
├── contexts/              # React context providers
│   └── AppContext.tsx     # Global application state management
└── hooks/                 # Custom React hooks
    ├── useAppContext.ts   # Hook for accessing app context
    ├── useAutocompleteSuggestions.ts # Autocomplete functionality
    └── useGeocoder.ts     # Google Maps geocoder hook
```

### Architecture Patterns

- **Context API**: Global state management via `AppContext` for location, entities, address, and view state
- **Custom Hooks**: Reusable logic extracted into custom hooks (`useAppContext`, `useGeocoder`, `useAutocompleteSuggestions`)
- **Components**: Modular components organized by feature (Sidebar, Map, Search)
- **Type Safety**: Centralized type definitions in `types.d.ts`
- **API Layer**: Separated API functions in `api.ts` for data fetching

### Application Flow

1. **Entry Point** (`main.tsx`): Wraps the app with `APIProvider` (Google Maps) and `AppProvider` (app state)
2. **Main App** (`App.tsx`): Renders the layout with Sidebar and Map components
3. **State Management**: `AppContext` manages location selection, entity fetching, and view switching
4. **User Interactions**: Users can search, click on the map, or use geolocation to set a location
5. **Data Fetching**: When a location is set, the app fetches entities and address information in parallel
6. **View Rendering**: The sidebar displays either an introduction view or details view based on the current state

## Styling

The application uses **Tailwind CSS v4** with custom design tokens defined in `src/styles.css`.
All design tokens are available as CSS custom properties and can be used throughout the application via Tailwind utilities.
