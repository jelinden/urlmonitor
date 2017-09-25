var urls;
var names;

function fillEndpoints() {
    document.getElementById("view").innerHTML = '';
    for (var i in urls) {
        document.getElementById("view").innerHTML += "<div class=\"item\"><h2>" + 
            urls[i] + 
            "</h2><canvas id=\"" +names[i]+"\" width=\"300\" height=\"180\"></canvas>" +
            "<div id=\"img" + i +"\"></div></div>";
        getValues(urls[i], names[i], valuesCallback);
        setInterval(getValues, 60000, urls[i], names[i], valuesCallback);
    }
}

function draw(divID, valueArray) {
    var i = 0;
    var ctx = document.getElementById(divID).getContext('2d');
    var newChart = new Chart(ctx, {
        type: 'line',
        data: {
            labels: valueArray.firstLoads.map(function(x){return i+=1;}),
            datasets: [{
                label: '',
                data: valueArray.firstLoads,
                backgroundColor: 'rgba(75, 192, 192, 0.3)',
                borderColor: 'rgba(75, 192, 192, 1)',
                borderWidth: 1
            },
            {
                label: '',
                data: valueArray.renderLoads,
                backgroundColor: 'rgba(255, 200, 132, 0.1)',
                borderColor: 'rgba(255,200,132, 1)',
                borderWidth: 1
            }]
        },
        options: {
            scales: {
                yAxes: [{
                    ticks: {
                        beginAtZero:true,
                        steps: 10,
                        stepValue: 1,
                        max: 8
                    }
                }]
            },
            legend: {
                display: false
            }
        }
    });
}

function getURLs() {
    var xhr = new XMLHttpRequest();
    xhr.onload = function() {
        urls = JSON.parse(this.responseText).urls
        names = JSON.parse(this.responseText).names
        fillEndpoints();
    };

    xhr.open('GET', '/urls', true);
    xhr.send();
}

function getValues(url, divID, callback) {
    request('/json/' + encodeURIComponent(url), callback, divID);
}

function valuesCallback(resp, divID) {
    var parsed = JSON.parse(resp);
    var times = parsed.times;
    var div = document.getElementById(divID);
    var nextSibling = div.nextSibling;
    
    nextSibling.innerHTML = "<div>"+parsed.times.latest+"</div>";
    addimage(nextSibling, "/assets/img/" + parsed.times.service + ".png");
    draw(divID, times);
}

function addimage(sibbling, url) {
    var img = new Image();
    img.src = url+'?'+(new Date()).getMinutes();
    img.width = "300";
    sibbling.appendChild(img);
}

function request(url, callback, divID) {
    var xmlhttp = new XMLHttpRequest();
    xmlhttp.onload = function() {
        callback(xmlhttp.responseText, divID);
    };
    xmlhttp.open("GET", url, true);
    xmlhttp.send();
}
  
window.addEventListener("load", function(event) {
    getURLs();
});