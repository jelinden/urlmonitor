var urls;


function fillEndpoints() {
    document.getElementById("view").innerHTML = '';
    for (var i in urls) {
        document.getElementById("view").innerHTML += "<div class=\"item\"><h2>" + 
            urls[i] + 
            "</h2><canvas id=\"chart" + i + "\" width=\"300\" height=\"180\"></canvas></div>";
        getValues(urls[i], "chart" + i, valuesCallback);
        setInterval(getValues, 30000, urls[i], "chart" + i, valuesCallback);
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
                data: valueArray.finalLoads,
                backgroundColor: 'rgba(255, 99, 132, 0.1)',
                borderColor: 'rgba(255,99,132,1)',
                borderWidth: 1
            }]
        },
        options: {
            scales: {
                yAxes: [{
                    ticks: {
                        beginAtZero:true
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
        fillEndpoints();
    };

    xhr.open('GET', '/urls', true);
    xhr.send();
}

function getValues(url, divID, callback) {
    request('/json/' + encodeURIComponent(url), callback, divID);
}

function valuesCallback(resp, divID) {
    var times = JSON.parse(resp).times;
    draw(divID, times);
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