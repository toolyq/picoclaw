package commands

import "context"

func resetCommand() Definition {
	return Definition{
		Name:        "reset",
		Description: "Clear current session history and memory",
		Handler: func(_ context.Context, req Request, rt *Runtime) error {
			if rt == nil || rt.ClearSession == nil {
				return req.Reply(unavailableMsg)
			}
			if err := rt.ClearSession(); err != nil {
				return req.Reply("Failed to clear session: " + err.Error())
			}
			return req.Reply("Session cleared.")
		},
	}
}
