package server
import("encoding/json";"net/http";"strconv";"github.com/stockyard-dev/stockyard-concierge/internal/store")
func(s *Server)handleList(w http.ResponseWriter,r *http.Request){list,_:=s.db.List();if list==nil{list=[]store.Checklist{}};writeJSON(w,200,list)}
func(s *Server)handleCreate(w http.ResponseWriter,r *http.Request){var c store.Checklist;json.NewDecoder(r.Body).Decode(&c);if c.Name==""{writeError(w,400,"name required");return};s.db.Create(&c);writeJSON(w,201,c)}
func(s *Server)handleDelete(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);s.db.Delete(id);writeJSON(w,200,map[string]string{"status":"deleted"})}
func(s *Server)handleGetProgress(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);userID:=r.URL.Query().Get("user_id");if userID==""{writeError(w,400,"user_id required");return};p,_:=s.db.GetProgress(id,userID);writeJSON(w,200,p)}
func(s *Server)handleUpdateProgress(w http.ResponseWriter,r *http.Request){id,_:=strconv.ParseInt(r.PathValue("id"),10,64);var req struct{UserID string `json:"user_id"`;CompletedSteps string `json:"completed_steps"`;Percent int `json:"percent"`};json.NewDecoder(r.Body).Decode(&req);if req.UserID==""{writeError(w,400,"user_id required");return};s.db.UpdateProgress(id,req.UserID,req.CompletedSteps,req.Percent);writeJSON(w,200,map[string]string{"status":"updated"})}
func(s *Server)handleOverview(w http.ResponseWriter,r *http.Request){m,_:=s.db.Stats();writeJSON(w,200,m)}
