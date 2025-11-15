# Orchestrator Server (City Day Navigator)

This is the brain of the City Day Navigator. It is a Go server that connects the client app to the `mcp-server` (the "tools").

Its sole responsibility is to:
1. Receive a simple user request (e.g., "Plan my day in Kyoto").
2. Use the Gemini API for multi pass reasoning.
3. Call the `mcp-server` tools to gather live data (weather, locations, ETAs).
4. Synthesize all data into a final, narrative plan.
5. Stream the plan and a tool trace back to the client.

---

## 🚀 Setup & Run

1. **Clone the repo:**
   ```bash
   git clone https://github.com/ayushh2k/city-nav-orchestrator
   cd city-nav-orchestrator
   ```

2. **Install dependencies:**
   ```bash
   go mod tidy
   make install-air
   ```

3. **Create `.env` file:**
   ```
   # Get from Google AI Studio
   GEMINI_API_KEY="your_gemini_key"
   
   # URL of your deployed mcp-server
   MCP_SERVER_BASE_URL="http://127.0.0.1:8000" 
   MCP_SERVER_API_KEY="my-secret-key-123"
   ```

4. **Run the dev server:**
   ```bash
   make dev
   ```
   The server will run on `http://localhost:8080`.