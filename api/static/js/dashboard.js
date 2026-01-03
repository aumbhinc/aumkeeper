// Sidebar toggle
document.getElementById('menu-btn').addEventListener('click', ()=>{
  document.querySelector('aside').classList.toggle('collapsed');
});

// Theme toggle
document.querySelectorAll('.theme-toggler span').forEach(t=>{
  t.addEventListener('click', ()=>{
    document.body.classList.toggle('dark-theme');
    document.querySelectorAll('.theme-toggler span').forEach(s=>s.classList.toggle('active'));
  });
});

// Collapsible panels
document.querySelectorAll('.right .collapsible').forEach(c=>{
  c.addEventListener('click', ()=>{
    const content = c.nextElementSibling;
    content.style.display = (content.style.display === 'flex') ? 'none' : 'flex';
    c.querySelector('span').textContent = (content.style.display === 'flex') ? 'expand_less' : 'expand_more';
  });
});

// Charts
new Chart(document.getElementById('salesChart').getContext('2d'), {
  type: 'line',
  data: { labels:['Mon','Tue','Wed','Thu','Fri','Sat','Sun'], datasets:[{label:'Sales ($)', data:[12000,15000,13000,17000,18000,19000,22000], borderColor:'#00ffd8', backgroundColor:'rgba(0,255,216,0.2)', tension:0.3, fill:true, pointRadius:4, pointBackgroundColor:'#00ffd8'}]},
  options: { responsive:true, plugins:{legend:{labels:{color:'#e4e6eb'}}}, scales:{x:{ticks:{color:'#e4e6eb'}}, y:{ticks:{color:'#e4e6eb'}}} }
});

new Chart(document.getElementById('accountsChart').getContext('2d'), {
  type:'doughnut',
  data:{ labels:['Assets','Liabilities','Equity','Revenue','Expenses'], datasets:[{data:[40,20,15,15,10], backgroundColor:['#34c759','#ff9f0a','#00ffd8','#ff453a','#9ca3af']}] },
  options:{ responsive:true, plugins:{legend:{position:'bottom', labels:{color:'#e4e6eb'}}} }
});
