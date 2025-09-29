Next Step: Implement "Goal Detail & History" Feature
Backend Tasks:

Create Endpoint to Get a Single Goal:

Endpoint: GET /api/v1/goals/:id

Function: Accepts a goal id as a URL parameter and returns the full details for that single goal object.

Create Endpoint to Get Check-in History:

Endpoint: GET /api/v1/checkins

Function: Accepts a goal_id as a query parameter and returns an array of all check-in records associated with that goal.

Frontend Tasks:

Modify Dashboard.vue:

Wrap each GoalCard component in a <router-link> to make it a clickable link.

The link's destination must be the dynamic route for the goal's detail page, e.g., :to="'/goals/' + goal.id".

Simplify GoalCard.vue Component:

Remove the "Completed", "Partial", and "Failed" buttons and their associated logic from this component. Its purpose is now for display and navigation.

Create New GoalDetail.vue View:

Routing: Create a new dynamic route at /goals/:id that maps to this component.

Logic:

On component mount, get the id from the URL route parameters.

Call the GET /api/v1/goals/:id API to fetch and display the goal's title.

Call the GET /api/v1/checkins?goal_id=:id API to fetch and display a list of all historical check-in records (showing date, status, and review notes).