/*
 * Task Management Full-Stack (React UI, no JSX transpilation required).
 * React 16 UMD build bundled locally for an installation-free Windows demo.
 * Core screens reflect the 2025 team's task fields & routes.
 */
(function () {
  'use strict';
  var h = React.createElement;
  function el(tag, props) { return h.apply(null, [tag, props].concat([].slice.call(arguments, 2))); }
  function str(s) { return String(s == null ? '' : s); }
  var statusNames = { todo: 'К выполнению', in_progress: 'В работе', done: 'Готово' };
  var priorityNames = { low: 'Низкий', medium: 'Средний', high: 'Высокий' };
  var initialForm = { title: '', description: '', due_date: '', priority: 'medium', status: 'todo' };
  function dateLabel(s) {
    if (!s) return 'Без срока';
    var d = s.split('-'); return d.length === 3 ? d[2] + '.' + d[1] + '.' + d[0] : s;
  }
  function today() { var d = new Date(); return d.getFullYear() + '-' + String(d.getMonth()+1).padStart(2,'0') + '-' + String(d.getDate()).padStart(2,'0'); }
  function icon(name) {
    var data = {
      grid: '<rect x="3" y="3" width="7" height="7" rx="1"/><rect x="14" y="3" width="7" height="7" rx="1"/><rect x="3" y="14" width="7" height="7" rx="1"/><rect x="14" y="14" width="7" height="7" rx="1"/>',
      list: '<rect x="9" y="4" width="12" height="3" rx="1"/><rect x="9" y="11" width="12" height="3" rx="1"/><rect x="9" y="18" width="12" height="3" rx="1"/><circle cx="4" cy="5.5" r="1"/><circle cx="4" cy="12.5" r="1"/><circle cx="4" cy="19.5" r="1"/>',
      calendar: '<rect x="3" y="5" width="18" height="16" rx="2"/><path d="M7 3v4M17 3v4M3 10h18"/>',
      plus: '<path d="M12 5v14M5 12h14"/>',
      download: '<path d="M12 3v12m-5-5 5 5 5-5M4 18v3h16v-3"/>',
      search: '<circle cx="11" cy="11" r="7"/><path d="m16 16 5 5"/>',
      arrow: '<path d="m8 5 7 7-7 7"/>',
      check: '<path d="m4 12 5 5L20 6"/>',
      trash: '<path d="M4 7h16m-10 4v6m4-6v6M6 7l1 14h10l1-14M9 7V4h6v3"/>',
      pencil: '<path d="m3 17 12-12 4 4-12 12-5 1 1-5Zm10-10 4 4"/>',
      flag: '<path d="M5 22V4m0 1q5-4 10 0t7 0v11q-4 4-9 0t-8 0"/>',
      clock: '<circle cx="12" cy="12" r="9"/><path d="M12 7v6l4 2"/>',
      close: '<path d="M6 6l12 12M18 6 6 18"/>',
      chart: '<path d="M4 20V11m6 9V5m6 15v-8m5 8H3"/>',
      spark: '<path d="m12 2 2.2 7.8L22 12l-7.8 2.2L12 22l-2.2-7.8L2 12l7.8-2.2z"/>',
      refresh: '<path d="M20 7a9 9 0 1 0 1 7M20 3v5h-5"/>'
    };
    return el('svg', { viewBox: '0 0 24 24', width: 18, height: 18, fill:'none', stroke:'currentColor', strokeWidth:1.7, strokeLinecap:'round', strokeLinejoin:'round', 'aria-hidden':'true', dangerouslySetInnerHTML:{__html:data[name]||data.grid} });
  }
  async function api(url,opts) {
    var options = Object.assign({headers:{'Content-Type':'application/json'}},opts||{});
    var r = await fetch(url,options);
    if(r.status === 204) return null;
    var b=await r.json().catch(function(){return {};});
    if(!r.ok) throw new Error(b.error||'HTTP '+r.status);
    return b;
  }
  function NavLink(props) {
    return el('a',{href:props.href,className:'nav-link'+(props.active?' active':''),onClick:props.click},icon(props.icon),el('span',null,props.text),props.count!=null?el('span',{className:'nav-count'},props.count):null);
  }
  function Tag(props) {return el('span',{className:'tag '+props.type+'-'+props.value},props.children||props.names[props.value]||props.value);}
  function Metric(props) {
    return el('div',{className:'metric'},
      el('div',{className:'metric-top'},el('span',null,props.title),el('span',{className:'metric-icon'},icon(props.icon))),
      el('div',{className:'metric-value'},str(props.value)),
      el('div',{className:'metric-foot'},props.hint));
  }
  function Empty(props) {
    return el('div',{className:'empty'},el('div',{className:'empty-image'},icon('list')),
      el('h3',null,'Пока здесь нет задач'),el('p',null,props.search?'Попробуй изменить поиск или фильтры.':'Создай первую задачу — она появится здесь и сохранится.'),
      !props.search?el('button',{className:'btn btn-primary',onClick:props.onNew},icon('plus'),' Новая задача'):null);
  }
  class App extends React.Component {
    constructor(props) {
      super(props);
      this.state={tasks:[],loaded:false,error:'',toast:'',page:location.pathname==='/board'?'board':location.pathname==='/stats'?'stats':'tasks',q:'',status:'all',priority:'all',modal:false,editId:null,form:Object.assign({},initialForm),saving:false,deleteId:null,sidebar:false};
      this.refresh=this.refresh.bind(this);
      this.onPop=this.onPop.bind(this);
    }
    componentDidMount(){this.refresh();window.addEventListener('popstate',this.onPop);if(/^\/tasks\/\d+\/edit$/.test(location.pathname)){var id=Number(location.pathname.split('/')[2]);this.editTaskByID(id);}if(location.pathname==='/tasks/new'){this.openNew();}}
    componentWillUnmount(){window.removeEventListener('popstate',this.onPop);}
    onPop(){this.setState({page:location.pathname==='/board'?'board':location.pathname==='/stats'?'stats':'tasks',modal:false,sidebar:false});}
    async refresh(){try{var tasks=await api('/api/tasks');this.setState({tasks:tasks,loaded:true,error:''});}catch(e){this.setState({loaded:true,error:'Не удалось подключиться к серверу: '+e.message});}}
    navigate(page){var url=page==='board'?'/board':page==='stats'?'/stats':'/tasks';history.pushState({},'',url);this.setState({page:page,modal:false,sidebar:false});}
    notify(msg){this.setState({toast:msg});}
    openNew(){this.setState({modal:true,editId:null,form:Object.assign({},initialForm)});history.pushState({},'', '/tasks/new');}
    editTaskByID(id){var task=this.state.tasks.find(function(t){return t.id===id;});if(!task){api('/api/tasks/'+id).then(function(t){this.openEdit(t);}.bind(this)).catch(function(){this.notify('Задача не найдена');}.bind(this));return;}this.openEdit(task);}
    openEdit(t){this.setState({modal:true,editId:t.id,form:{title:t.title,description:t.description||'',due_date:t.due_date||'',priority:t.priority||'medium',status:t.status||'todo'}});history.pushState({},'', '/tasks/'+t.id+'/edit');}
    closeModal(){this.setState({modal:false,editId:null});this.navigate(this.state.page);}
    async save(){var f=this.state.form;if(!f.title.trim()){this.notify('Введите название задачи');return;}this.setState({saving:true});try{if(this.state.editId){await api('/api/tasks/'+this.state.editId,{method:'PUT',body:JSON.stringify(f)});this.notify('Задача обновлена');}else{await api('/api/tasks',{method:'POST',body:JSON.stringify(f)});this.notify('Новая задача создана');}this.setState({saving:false});this.closeModal();await this.refresh();}catch(e){this.setState({saving:false});this.notify('Ошибка: '+e.message);}}
    async changeStatus(task,status){try{await api('/api/tasks/'+task.id,{method:'PUT',body:JSON.stringify(Object.assign({},task,{status:status}))});this.refresh();this.notify('Статус обновлён');}catch(e){this.notify('Ошибка: '+e.message);}}
    async remove(id){if(!window.confirm('Удалить задачу? Это действие нельзя отменить.'))return;try{await api('/api/tasks/'+id,{method:'DELETE'});this.refresh();this.notify('Задача удалена');}catch(e){this.notify('Ошибка: '+e.message);}}
    filtered(){var q=this.state.q.toLowerCase();return this.state.tasks.filter(function(t){return (!q || (t.title+' '+(t.description||'')).toLowerCase().includes(q)) && (this.state.status==='all'||this.state.status===t.status) && (this.state.priority==='all'||this.state.priority===t.priority);}.bind(this));}
    stats(){var ts=this.state.tasks;return {total:ts.length,done:ts.filter(function(t){return t.status==='done';}).length,progress:ts.filter(function(t){return t.status==='in_progress';}).length,overdue:ts.filter(function(t){return t.status!=='done'&&t.due_date&&t.due_date<today();}).length};}
    taskCard(t){return el('div',{className:'board-task',key:t.id,draggable:true,onDragStart:function(e){e.dataTransfer.setData('text/plain',t.id);}.bind(this)},
      el('div',{className:'card-header'},el('span',{className:'task-number'},'#'+String(t.id).padStart(3,'0')),el('button',{className:'icon-action',title:'Редактировать',onClick:function(){this.openEdit(t);}.bind(this)},icon('pencil'))),
      el('h4',null,t.title),el('p',null,t.description||'Без описания'),
      el('div',{className:'task-bottom'},el(Tag,{type:'priority',value:t.priority,names:priorityNames}),el('span',{className:'due'+(t.due_date&&t.due_date<today()&&t.status!=='done'?' overdue':'')},icon('calendar'),' ',dateLabel(t.due_date))));}
    board(){var cols=[['todo','К выполнению','status-todo'],['in_progress','В работе','status-in_progress'],['done','Готово','status-done']];var ts=this.filtered();return el('div',{className:'board-columns'},cols.map(function(item){var list=ts.filter(function(t){return t.status===item[0]});return el('section',{className:'board-column',key:item[0],onDragOver:function(e){e.preventDefault();},onDrop:function(e){e.preventDefault();var id=Number(e.dataTransfer.getData('text/plain'));var t=this.state.tasks.find(function(t){return t.id===id;});if(t&&t.status!==item[0])this.changeStatus(t,item[0]);}.bind(this)},
      el('div',{className:'column-head'},el('div',null,el('span',{className:'column-dot '+item[2]}),el('strong',null,item[1])),el('span',{className:'column-count'},list.length)),
      el('div',{className:'column-body'},list.map(this.taskCard.bind(this)),list.length===0?el('div',{className:'board-empty'},'Перетащи задачу сюда'):null));}.bind(this)));}
    row(t){return el('div',{className:'task-row',key:t.id},
      el('button',{title:t.status==='done'?'Вернуть в работу':'Завершить задачу',className:'check-btn'+(t.status==='done'?' checked':''),onClick:function(){this.changeStatus(t,t.status==='done'?'todo':'done');}.bind(this)},t.status==='done'?icon('check'):null),
      el('div',{className:'task-title-cell'},el('strong',{className:t.status==='done'?'crossed':''},t.title),el('small',null,t.description||'Без описания')),
      el('div',{className:'desktop-col'},el(Tag,{type:'status',value:t.status,names:statusNames})),
      el('div',{className:'desktop-col'},el(Tag,{type:'priority',value:t.priority,names:priorityNames})),
      el('div',{className:'due'+(t.due_date&&t.due_date<today()&&t.status!=='done'?' overdue':'')},dateLabel(t.due_date)),
      el('div',{className:'row-actions'},el('button',{className:'icon-action',title:'Изменить',onClick:function(){this.openEdit(t);}.bind(this)},icon('pencil')),el('button',{className:'icon-action danger',title:'Удалить',onClick:function(){this.remove(t.id);}.bind(this)},icon('trash'))));}
    table(){var ts=this.filtered();return el('section',{className:'tasks-panel'},el('div',{className:'panel-title'},el('h2',null,'Все задачи'),el('span',null,ts.length+' задач')),
      el('div',{className:'tasks-table-head'},el('span',null,'ЗАДАЧА'),el('span',null,'СТАТУС'),el('span',null,'ПРИОРИТЕТ'),el('span',null,'СРОК'),el('span',null,'ДЕЙСТВИЯ')),
      ts.length?el('div',null,ts.map(this.row.bind(this))):el(Empty,{onNew:this.openNew.bind(this),search:!!(this.state.q||this.state.status!=='all'||this.state.priority!=='all')}));}
    dashboard(){var s=this.stats(),percentage=s.total?Math.round(s.done/s.total*100):0;return el('div',null,
      el('div',{className:'metrics-grid'},el(Metric,{title:'Всего задач',value:s.total,icon:'list',hint:'Все созданные задачи'}),el(Metric,{title:'В работе',value:s.progress,icon:'clock',hint:'Задачи в процессе'}),el(Metric,{title:'Выполнено',value:s.done,icon:'check',hint:percentage+'% от общего числа'}),el(Metric,{title:'Просрочено',value:s.overdue,icon:'flag',hint:'Требуют внимания'})),
      el('section',{className:'progress-panel'},el('div',{className:'progress-head'},el('div',null,el('h3',null,'Твой прогресс'),el('p',null,'Маленькие шаги приводят к большому результату.')),el('div',{className:'progress-number'},percentage+'%')),
        el('div',{className:'progress-track'},el('div',{className:'progress-bar',style:{width:percentage+'%'}})),el('div',{className:'progress-caption'},s.done+' завершено из '+s.total+' задач')),
      this.table());}
    statsPage(){var s=this.stats();var groups=[['К выполнению',s.total-s.done-s.progress],['В работе',s.progress],['Готово',s.done],['Просрочено',s.overdue]];return el('div',{className:'stats-page'},el('div',{className:'metrics-grid'},el(Metric,{title:'Всего задач',value:s.total,icon:'list',hint:'Всего добавлено'}),el(Metric,{title:'В работе',value:s.progress,icon:'clock',hint:'Активные задачи'}),el(Metric,{title:'Готово',value:s.done,icon:'check',hint:'Завершены'}),el(Metric,{title:'Просрочено',value:s.overdue,icon:'flag',hint:'Проверь сроки'})),el('section',{className:'stats-panel'},el('h3',null,'Распределение задач'),groups.map(function(pair){var pct=s.total?Math.round(pair[1]/s.total*100):0;return el('div',{className:'stats-line',key:pair[0]},el('div',{className:'stats-label'},el('span',null,pair[0]),el('b',null,pair[1])),el('div',{className:'progress-track'},el('div',{className:'progress-bar',style:{width:pct+'%'}})));}),el('a',{href:'/api/export.csv',className:'btn btn-outline'},icon('download'),' Скачать CSV')));}
    formModal(){var f=this.state.form;var fields=[['title','Название задачи','text','Например, подготовить презентацию'],['description','Описание','textarea','Что нужно сделать?'],['due_date','Срок выполнения','date','']];var onChange=function(key,value){this.setState({form:Object.assign({},this.state.form,{[key]:value})});}.bind(this);
      return el('div',{className:'modal-backdrop',onMouseDown:function(e){if(e.target===e.currentTarget)this.closeModal();}.bind(this)},el('div',{className:'modal-dialog',role:'dialog','aria-modal':'true','aria-label':this.state.editId?'Редактировать задачу':'Новая задача'},
        el('div',{className:'modal-head'},el('div',null,el('small',null,'TASK MANAGEMENT'),el('h2',null,this.state.editId?'Редактировать задачу':'Новая задача')),el('button',{className:'icon-action',onClick:this.closeModal.bind(this),'aria-label':'Закрыть'},icon('close'))),
        el('form',{onSubmit:function(e){e.preventDefault();this.save();}.bind(this)},
          fields.map(function(field){return el('label',{className:'form-field',key:field[0]},el('span',null,field[1],field[0]==='title'?' *':''),field[2]==='textarea'?el('textarea',{value:f[field[0]],placeholder:field[3],rows:4,maxLength:2000,onChange:function(e){onChange(field[0],e.target.value);}}):el('input',{type:field[2],value:f[field[0]],placeholder:field[3],required:field[0]==='title',maxLength:field[0]==='title'?160:undefined,onChange:function(e){onChange(field[0],e.target.value);}}));}),
          el('div',{className:'form-pair'},el('label',{className:'form-field'},el('span',null,'Приоритет'),el('select',{value:f.priority,onChange:function(e){onChange('priority',e.target.value);}},Object.keys(priorityNames).map(function(k){return el('option',{value:k,key:k},priorityNames[k]);}))),el('label',{className:'form-field'},el('span',null,'Статус'),el('select',{value:f.status,onChange:function(e){onChange('status',e.target.value);}},Object.keys(statusNames).map(function(k){return el('option',{value:k,key:k},statusNames[k]);})))),
          el('div',{className:'modal-actions'},el('button',{type:'button',className:'btn btn-outline',onClick:this.closeModal.bind(this)},'Отмена'),el('button',{type:'submit',className:'btn btn-primary',disabled:this.state.saving},this.state.saving?'Сохраняю...':(this.state.editId?'Сохранить изменения':'Создать задачу'))))));}
    async loadExamples(){var entries=[{title:'Подготовить финальный отчёт',description:'Собрать материалы и проверить результаты.',due_date:'',priority:'high',status:'in_progress'},{title:'Сделать презентацию проекта',description:'10 слайдов: цель, стек, результат, demo.',due_date:'',priority:'medium',status:'todo'},{title:'Настроить GitHub репозиторий',description:'Добавить README и инструкции запуска.',due_date:'',priority:'low',status:'done'}];try{for(var i=0;i<entries.length;i++)await api('/api/tasks',{method:'POST',body:JSON.stringify(entries[i])});this.notify('Примеры задач добавлены');this.refresh();}catch(e){this.notify('Ошибка: '+e.message);}}
    render(){var s=this.stats(),page=this.state.page;return el('div',{className:'app-shell'},
      el('aside',{className:'sidebar'+(this.state.sidebar?' open':'')},
        el('div',{className:'brand'},el('div',{className:'brand-icon'},icon('grid')),el('div',null,el('strong',null,'taskflow'),el('small',null,'PROJECT WORKSPACE'))),
        el('div',{className:'menu-caption'},'РАБОЧЕЕ ПРОСТРАНСТВО'),
        el(NavLink,{href:'/tasks',active:page==='tasks',icon:'grid',text:'Обзор',click:function(e){e.preventDefault();this.navigate('tasks');}.bind(this)}),
        el(NavLink,{href:'/board',active:page==='board',icon:'list',text:'Kanban-доска',click:function(e){e.preventDefault();this.navigate('board');}.bind(this)}),
        el(NavLink,{href:'/stats',active:page==='stats',icon:'chart',text:'Статистика',click:function(e){e.preventDefault();this.navigate('stats');}.bind(this)}),
        el('div',{className:'sidebar-divider'}),el('div',{className:'menu-caption'},'БЫСТРЫЙ ДОСТУП'),
        el('button',{className:'nav-link side-action',onClick:function(){this.setState({status:'done'});this.navigate('tasks');}.bind(this)},icon('check'),el('span',null,'Завершённые'),el('span',{className:'nav-count'},s.done)),
        el('a',{href:'/api/export.csv',className:'nav-link'},icon('download'),el('span',null,'Экспорт CSV')),
        el('div',{className:'sidebar-spacer'}),
        el('div',{className:'sidebar-help'},el('div',{className:'help-symbol'},icon('spark')),el('h4',null,'Порядок в задачах'),el('p',null,'Планируй день и отслеживай прогресс в одном месте.'),el('button',{onClick:this.openNew.bind(this)},'Создать задачу ',icon('arrow'))),
        el('div',{className:'sidebar-footer'},el('div',{className:'user-avatar'},'Z'),el('div',null,el('strong',null,'Мой проект'),el('span',null,'Local workspace')))),
      el('main',{className:'main'},
        el('header',{className:'topbar'},el('div',{className:'top-left'},el('button',{className:'menu-toggle',onClick:function(){this.setState({sidebar:!this.state.sidebar});}.bind(this)},icon('list')),el('span',null,'Workspace'),el('span',{className:'top-separator'},'/'),el('strong',null,page==='board'?'Kanban-доска':page==='stats'?'Статистика':'Мои задачи')),
          el('div',{className:'top-right'},el('span',{className:'online-dot'}),'Локальный проект',el('span',{className:'date-chip'},dateLabel(today())))),
        el('div',{className:'content'},
          el('section',{className:'hero'},el('div',null,el('div',{className:'eyebrow'},'TASK MANAGEMENT · FULL STACK 2025'),el('h1',null,page==='board'?'Доска задач':page==='stats'?'Аналитика задач':'Управляй своим временем',el('span',{className:'title-period'},'.')),el('p',null,page==='board'?'Перемещай задачи между этапами и отслеживай выполнение.':page==='stats'?'Обзор эффективности и прогресса по задачам.':'Здесь собраны все твои задачи. Сфокусируйся на главном и двигайся вперёд.')),
            el('button',{className:'btn btn-primary',onClick:this.openNew.bind(this)},icon('plus'),' Новая задача')),
          this.state.error?el('div',{className:'error-banner'},this.state.error,el('button',{onClick:this.refresh},'Повторить')):null,
          el('div',{className:'toolbar'},el('label',{className:'search-box'},icon('search'),el('input',{value:this.state.q,placeholder:'Поиск по задачам...',onChange:function(e){this.setState({q:e.target.value});}.bind(this),'aria-label':'Поиск задач'})),
            el('div',{className:'filter-group'},el('select',{value:this.state.status,onChange:function(e){this.setState({status:e.target.value});}.bind(this),'aria-label':'Фильтр статуса'},el('option',{value:'all'},'Все статусы'),Object.keys(statusNames).map(function(k){return el('option',{value:k,key:k},statusNames[k]);})),el('select',{value:this.state.priority,onChange:function(e){this.setState({priority:e.target.value});}.bind(this),'aria-label':'Фильтр приоритета'},el('option',{value:'all'},'Все приоритеты'),Object.keys(priorityNames).map(function(k){return el('option',{value:k,key:k},priorityNames[k]);})),el('button',{className:'icon-action refresh-btn',title:'Обновить',onClick:this.refresh},icon('refresh')))),
          !this.state.loaded?el('div',{className:'loading'},'Загружаю задачи...'):page==='board'?this.board():page==='stats'?this.statsPage():this.dashboard(),
          this.state.loaded&&s.total===0&&!this.state.q&&page==='tasks'?el('div',{className:'demo-hint'},'Хочешь увидеть работу приложения? ',el('button',{onClick:this.loadExamples.bind(this)},'Добавить 3 демонстрационные задачи')):null,
          el('footer',{className:'footer'},'Task Management Full-Stack · Go + React · 2025 / Restored Edition',el('span',null,'Local API: /api/tasks')))),
      this.state.modal?this.formModal():null,
      this.state.toast?el('div',{className:'toast',role:'status',onClick:function(){this.setState({toast:''});}.bind(this)},this.state.toast,' ×'):null);
    }
  }
  ReactDOM.render(el(App,null),document.getElementById('root'));
})();
