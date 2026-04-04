package server

import "net/http"

func (s *Server) dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(dashHTML))
}

const dashHTML = `<!DOCTYPE html><html><head><meta charset="UTF-8"><meta name="viewport" content="width=device-width,initial-scale=1.0"><title>Concierge</title>
<link href="https://fonts.googleapis.com/css2?family=JetBrains+Mono:wght@400;500;700&display=swap" rel="stylesheet">
<style>
:root{--bg:#1a1410;--bg2:#241e18;--bg3:#2e261e;--rust:#e8753a;--leather:#a0845c;--cream:#f0e6d3;--cd:#bfb5a3;--cm:#7a7060;--gold:#d4a843;--green:#4a9e5c;--red:#c94444;--blue:#5b8dd9;--mono:'JetBrains Mono',monospace}
*{margin:0;padding:0;box-sizing:border-box}body{background:var(--bg);color:var(--cream);font-family:var(--mono);line-height:1.5}
.hdr{padding:1rem 1.5rem;border-bottom:1px solid var(--bg3);display:flex;justify-content:space-between;align-items:center}.hdr h1{font-size:.9rem;letter-spacing:2px}.hdr h1 span{color:var(--rust)}
.main{padding:1.5rem;max-width:960px;margin:0 auto}
.stats{display:grid;grid-template-columns:repeat(3,1fr);gap:.5rem;margin-bottom:1rem}
.st{background:var(--bg2);border:1px solid var(--bg3);padding:.6rem;text-align:center}
.st-v{font-size:1.2rem;font-weight:700}.st-l{font-size:.5rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-top:.15rem}
.toolbar{display:flex;gap:.5rem;margin-bottom:1rem;align-items:center}
.search{flex:1;padding:.4rem .6rem;background:var(--bg2);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.search:focus{outline:none;border-color:var(--leather)}
.item{background:var(--bg2);border:1px solid var(--bg3);padding:.8rem 1rem;margin-bottom:.5rem;transition:border-color .2s}
.item:hover{border-color:var(--leather)}
.item-top{display:flex;justify-content:space-between;align-items:flex-start;gap:.5rem}
.item-title{font-size:.85rem;font-weight:700}
.item-sub{font-size:.7rem;color:var(--cd);margin-top:.1rem}
.progress-bar{height:6px;background:var(--bg3);margin-top:.4rem;border-radius:3px;overflow:hidden}
.progress-fill{height:100%;background:var(--green);transition:width .3s}
.item-meta{font-size:.55rem;color:var(--cm);margin-top:.3rem;display:flex;gap:.5rem;flex-wrap:wrap;align-items:center}
.item-actions{display:flex;gap:.3rem;flex-shrink:0}
.badge{font-size:.5rem;padding:.12rem .35rem;text-transform:uppercase;letter-spacing:1px;border:1px solid}
.badge.pending{border-color:var(--gold);color:var(--gold)}.badge.in-progress{border-color:var(--blue);color:var(--blue)}.badge.complete{border-color:var(--green);color:var(--green)}
.btn{font-size:.6rem;padding:.25rem .5rem;cursor:pointer;border:1px solid var(--bg3);background:var(--bg);color:var(--cd);transition:all .2s}
.btn:hover{border-color:var(--leather);color:var(--cream)}.btn-p{background:var(--rust);border-color:var(--rust);color:#fff}
.btn-sm{font-size:.55rem;padding:.2rem .4rem}
.modal-bg{display:none;position:fixed;inset:0;background:rgba(0,0,0,.65);z-index:100;align-items:center;justify-content:center}.modal-bg.open{display:flex}
.modal{background:var(--bg2);border:1px solid var(--bg3);padding:1.5rem;width:460px;max-width:92vw;max-height:90vh;overflow-y:auto}
.modal h2{font-size:.8rem;margin-bottom:1rem;color:var(--rust);letter-spacing:1px}
.fr{margin-bottom:.6rem}.fr label{display:block;font-size:.55rem;color:var(--cm);text-transform:uppercase;letter-spacing:1px;margin-bottom:.2rem}
.fr input,.fr select,.fr textarea{width:100%;padding:.4rem .5rem;background:var(--bg);border:1px solid var(--bg3);color:var(--cream);font-family:var(--mono);font-size:.7rem}
.fr input:focus,.fr select:focus,.fr textarea:focus{outline:none;border-color:var(--leather)}
.row2{display:grid;grid-template-columns:1fr 1fr;gap:.5rem}
.acts{display:flex;gap:.4rem;justify-content:flex-end;margin-top:1rem}
.empty{text-align:center;padding:3rem;color:var(--cm);font-style:italic;font-size:.75rem}
</style></head><body>
<div class="hdr"><h1><span>&#9670;</span> CONCIERGE</h1><button class="btn btn-p" onclick="openForm()">+ New Checklist</button></div>
<div class="main">
<div class="stats" id="stats"></div>
<div class="toolbar"><input class="search" id="search" placeholder="Search checklists..." oninput="render()"></div>
<div id="list"></div>
</div>
<div class="modal-bg" id="mbg" onclick="if(event.target===this)closeModal()"><div class="modal" id="mdl"></div></div>
<script>
var A='/api',items=[],editId=null;
async function load(){var r=await fetch(A+'/checklists').then(function(r){return r.json()});items=r.checklists||[];renderStats();render();}
function renderStats(){var total=items.length;var done=items.filter(function(c){return c.status==='complete'}).length;var active=items.filter(function(c){return c.status==='in-progress'}).length;
document.getElementById('stats').innerHTML='<div class="st"><div class="st-v">'+total+'</div><div class="st-l">Total</div></div><div class="st"><div class="st-v" style="color:var(--blue)">'+active+'</div><div class="st-l">Active</div></div><div class="st"><div class="st-v" style="color:var(--green)">'+done+'</div><div class="st-l">Complete</div></div>';}
function render(){var q=(document.getElementById('search').value||'').toLowerCase();var f=items;
if(q)f=f.filter(function(c){return(c.customer_name||'').toLowerCase().includes(q)||(c.template||'').toLowerCase().includes(q)||(c.assignee||'').toLowerCase().includes(q)});
if(!f.length){document.getElementById('list').innerHTML='<div class="empty">No checklists.</div>';return;}
var h='';f.forEach(function(c){
var pct=c.progress||0;
h+='<div class="item"><div class="item-top"><div style="flex:1">';
h+='<div class="item-title">'+esc(c.customer_name)+'</div>';
if(c.template)h+='<div class="item-sub">'+esc(c.template)+'</div>';
h+='</div><div class="item-actions">';
h+='<button class="btn btn-sm" onclick="openEdit(''+c.id+'')">Edit</button>';
h+='<button class="btn btn-sm" onclick="del(''+c.id+'')" style="color:var(--red)">&#10005;</button>';
h+='</div></div>';
h+='<div class="progress-bar"><div class="progress-fill" style="width:'+pct+'%"></div></div>';
h+='<div class="item-meta">';
var st=(c.status||'pending').replace(/_/g,'-');
h+='<span class="badge '+st+'">'+st+'</span>';
h+='<span>'+pct+'% complete</span>';
if(c.assignee)h+='<span>'+esc(c.assignee)+'</span>';
if(c.due_date)h+='<span>Due: '+c.due_date+'</span>';
h+='</div></div>';});
document.getElementById('list').innerHTML=h;}
async function del(id){if(!confirm('Delete?'))return;await fetch(A+'/checklists/'+id,{method:'DELETE'});load();}
function formHTML(cl){var i=cl||{customer_name:'',template:'',assignee:'',status:'pending',due_date:'',progress:0};var isEdit=!!cl;
var h='<h2>'+(isEdit?'EDIT':'NEW')+' CHECKLIST</h2>';
h+='<div class="fr"><label>Customer *</label><input id="f-name" value="'+esc(i.customer_name)+'"></div>';
h+='<div class="row2"><div class="fr"><label>Template</label><input id="f-tmpl" value="'+esc(i.template)+'" placeholder="e.g. Enterprise Onboarding"></div>';
h+='<div class="fr"><label>Assignee</label><input id="f-assign" value="'+esc(i.assignee)+'"></div></div>';
h+='<div class="row2"><div class="fr"><label>Status</label><select id="f-status">';
['pending','in-progress','complete'].forEach(function(s){h+='<option value="'+s+'"'+(i.status===s?' selected':'')+'>'+s+'</option>';});
h+='</select></div><div class="fr"><label>Due Date</label><input id="f-due" type="date" value="'+esc(i.due_date)+'"></div></div>';
h+='<div class="fr"><label>Progress (%)</label><input id="f-prog" type="number" min="0" max="100" value="'+(i.progress||0)+'"></div>';
h+='<div class="acts"><button class="btn" onclick="closeModal()">Cancel</button><button class="btn btn-p" onclick="submit()">'+(isEdit?'Save':'Create')+'</button></div>';
return h;}
function openForm(){editId=null;document.getElementById('mdl').innerHTML=formHTML();document.getElementById('mbg').classList.add('open');}
function openEdit(id){var c=null;for(var j=0;j<items.length;j++){if(items[j].id===id){c=items[j];break;}}if(!c)return;editId=id;document.getElementById('mdl').innerHTML=formHTML(c);document.getElementById('mbg').classList.add('open');}
function closeModal(){document.getElementById('mbg').classList.remove('open');editId=null;}
async function submit(){var name=document.getElementById('f-name').value.trim();if(!name){alert('Customer required');return;}
var body={customer_name:name,template:document.getElementById('f-tmpl').value.trim(),assignee:document.getElementById('f-assign').value.trim(),status:document.getElementById('f-status').value,due_date:document.getElementById('f-due').value,progress:parseInt(document.getElementById('f-prog').value)||0};
if(editId){await fetch(A+'/checklists/'+editId,{method:'PUT',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
else{await fetch(A+'/checklists',{method:'POST',headers:{'Content-Type':'application/json'},body:JSON.stringify(body)});}
closeModal();load();}
function esc(s){if(!s)return'';var d=document.createElement('div');d.textContent=s;return d.innerHTML;}
document.addEventListener('keydown',function(e){if(e.key==='Escape')closeModal();});
load();
</script></body></html>`
