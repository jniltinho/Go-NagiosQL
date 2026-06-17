> **DEFERRED**: This capability is out of scope for the current change (`migrar-nagiosql-golang-echo`).
> The backend API must be complete and fully tested before the frontend is started.
> Frontend implementation will be tracked in a separate proposal: `migrar-nagiosql-frontend-vuejs`.
> All requirements below are preserved for that proposal.

---

## ADDED Requirements

### Requirement: Vue.js 3 + TypeScript SPA with Vue Router
The frontend SHALL be a Vue.js 3 Single Page Application written entirely in TypeScript, using Vue Router in history mode. Routes SHALL mirror NagiosQL's page structure: `/admin/monitoring`, `/admin/hosts`, `/admin/services`, `/admin/verify`, etc. The SPA SHALL be built with Vite. All component props, API response types, and store state SHALL be fully typed.

#### Scenario: SPA navigation
- **WHEN** a user navigates between pages in the application
- **THEN** the URL changes without a full page reload

#### Scenario: Direct URL access
- **WHEN** a user navigates directly to `/admin/hosts` in the browser
- **THEN** the Echo server serves `index.html` and Vue Router renders the correct view

### Requirement: Pinia for state management
The frontend SHALL use Pinia as the state management library. Stores SHALL be defined for: `useAuthStore` (current user session), `useDomainStore` (active data domain), and `useSettingsStore` (NagiosQL global settings). Each store SHALL be TypeScript-typed and use the Composition API setup syntax.

#### Scenario: Auth store persists user across page navigations
- **WHEN** a user navigates between views after login
- **THEN** the `useAuthStore` provides the current user without re-fetching from the API on each view

#### Scenario: Domain store filters API requests
- **WHEN** the user switches the active domain via the domain selector
- **THEN** the `useDomainStore` updates and all list views automatically reload with the new `domain_id`

### Requirement: Tailwind CSS v4 styling matching NagiosQL layout
The frontend SHALL use Tailwind CSS v4 for styling. The layout SHALL be visually compatible with NagiosQL: top navigation bar, left sidebar menu, main content area with breadcrumbs, and list/form views.

#### Scenario: Navigation sidebar
- **WHEN** the user is logged in
- **THEN** a sidebar is rendered with menu sections matching NagiosQL: Monitoring Objects, Configuration, Administration

#### Scenario: Responsive layout
- **WHEN** the browser window is narrower than 768px
- **THEN** the sidebar collapses to a hamburger menu

### Requirement: Lucide Icons
All icons in the UI SHALL use Lucide Icons (`lucide-vue-next`). No other icon library SHALL be used. Icons SHALL be tree-shaken — only imported icons are bundled.

#### Scenario: Icons render correctly
- **WHEN** a list page is loaded
- **THEN** action icons (edit, delete, copy, view) render as Lucide SVG icons

### Requirement: List views with sorting and pagination
All object list views SHALL support server-side pagination (page size configurable per user), column sorting (ascending/descending), and a search/filter input. This matches NagiosQL's list view behavior.

#### Scenario: Paginated host list
- **WHEN** a user opens `/admin/hosts` with 50+ hosts in the database
- **THEN** hosts are displayed paginated with navigation controls and a count of total records

#### Scenario: Column sort
- **WHEN** a user clicks a column header in the host list
- **THEN** the list re-sorts by that column and the sort direction is toggled

### Requirement: Form views with validation
All object form views SHALL perform client-side validation before submitting to the API. Required fields SHALL be highlighted. Error responses from the API SHALL be displayed inline.

#### Scenario: Required field validation
- **WHEN** a user submits a host form without `host_name`
- **THEN** the field is highlighted in red and submission is blocked

#### Scenario: API error display
- **WHEN** the API returns a 400 error with a field-specific message
- **THEN** the error is displayed next to the relevant field

### Requirement: Multi-select widgets for relation fields
Fields that reference multiple related objects (e.g., host contact groups, service templates, notification options) SHALL use a dual-list multi-select widget (available → selected), matching NagiosQL's YUI multiselect behavior.

#### Scenario: Multi-select contact groups
- **WHEN** editing a host, the user opens the contact_groups field
- **THEN** a dual-list widget shows available contact groups on the left and selected ones on the right

### Requirement: Monitoring dashboard
The `/admin/monitoring` page SHALL display counts of active and inactive objects per type (hosts, services, host groups, service groups, host templates, service templates), equivalent to NagiosQL's `monitoring.php`.

#### Scenario: Dashboard object counts
- **WHEN** a user visits `/admin/monitoring`
- **THEN** a summary table shows active/inactive counts for each object type in the current domain

### Requirement: Embedded in binary — no external file access
The compiled frontend dist SHALL be served entirely from the embedded filesystem. The production build SHALL NOT require any files on the host filesystem.

#### Scenario: Frontend loads without dist directory
- **WHEN** the `nagiosql` binary runs in a directory without a `frontend/dist` folder
- **THEN** the frontend still loads correctly from the embedded assets
