
var ws;
var last;
var radius = 25;
var canvas;
var ctx;
var BB;
var offsetX;
var offsetY;
var WIDTH;
var HEIGHT;
var dragok;
var startX;
var startY;
var shapes={};
var background = new Image();
background.src = "https://www.roomsketcher.com/wp-content/uploads/2017/11/RoomSketcher-Office-Floor-Plan-PID3529710-2D-bw-with-Labels.jpg";
var printer;

$(document).ready(function() {
	
	document.getElementById("user").textContent= getCookie("user");
	var output = document.getElementById("output");
	printer = function(message) {
		var d = document.createElement("div");
		d.textContent = message;
		output.appendChild(d);
	}
	

	
	// get canvas related references
	last = {};
	radius = 25;
	canvas=document.getElementById("canvas");
	ctx=canvas.getContext("2d");
	BB=canvas.getBoundingClientRect();
	offsetX=BB.left;
	offsetY=BB.top;
	WIDTH = canvas.width;
	HEIGHT = canvas.height;

	// drag related variables
	dragok = false;

	shapes={};
	//shapes.push({id: 0, x:10,y:100,width:30,height:30,fill:"#444444",isDragging:false});
	//shapes.push({id: 1, x:80,y:100,width:30,height:30,fill:"#ff550d",isDragging:false});
	//shapes.push({id: 2, x:150,y:100,text: "hello", r:50,fill:"#800080",isDragging:false});
	//shapes["3"] = {id:"3", x:200, text: "hey", y:100,r:40,fill:"#0c64e8",isDragging:false, text: "hello"};

	// listen for mouse events
	canvas.onmousedown = myDown;
	canvas.onmouseup = myUp;
	canvas.onmousemove = myMove;

	// Canvas background image
	background.onload = function(){
		ctx.drawImage(background,-130,-60, 1600, 950);   
	}


	open();
	
	sleep(500).then(() => {
		drawsend();
	});
	// call to draw the scene
	drawsend();
	
	document.getElementById("setBackground").onclick = function(evt) {
		var image = document.getElementById("background").value; 
		background.src = image;
		return false;
	};
	
	document.getElementById("join").onclick = function(evt) {
		var name = getUser();
		console.log("Joining as ", name);
		setCookie("user", name, 5);
		join(name);
		document.getElementById("user").textContent= getCookie("user");
		return false;
	};

	document.getElementById("delete").onclick = function(evt) {
		var user = getUser();
		delete shapes[user];
		drawsend();
		return false;
	}

	document.getElementById("setTopic").onclick = function(evt) {
		var user = getUser();
		if (user === "") {
			return
		}
		topic = document.getElementById("topic").value;
		link = document.getElementById("link").value;
		data = shapes[user];
		data["topic"] = topic;
		data["link"] = link;
		shapes[user] = data;
		drawsend();
	}

	document.getElementById("clear").onclick = function(evt) {
		var user = getUser();
		if (user === "") {
			return
		}
		delete shapes[user][topic];
		delete shapes[user][link];

	}
});

function getUser() {
	var user = document.getElementById("name").value; 
	if (!user) {
		user = getCookie("user");
	}
	if (!user) {
		alert("please input a username");
	}
	return user
}

function setCookie(cname, cvalue, exdays) {
  var d = new Date();
  d.setTime(d.getTime() + (exdays*24*60*60*1000));
  var expires = "expires="+ d.toUTCString();
  document.cookie = cname + "=" + cvalue + ";" + expires + ";path=/";
}

function getCookie(cname) {
  var name = cname + "=";
  var decodedCookie = decodeURIComponent(document.cookie);
  var ca = decodedCookie.split(';');
  for(var i = 0; i <ca.length; i++) {
    var c = ca[i];
    while (c.charAt(0) == ' ') {
      c = c.substring(1);
    }
    if (c.indexOf(name) == 0) {
      return c.substring(name.length, c.length);
    }
  }
  return "";
}


// draw a single rect
function rect(r) {

	const newDiv = document.createElement("div");
	newDiv.style.top = r.y;
	newDiv.style.left = r.x;
	newDiv.style.width = r.width;
	newDiv.style.height = r.height;
	newDiv.classList.add("overlay");
	const ref= document.createElement("a");
	if (r.link) {
		ref.href = r.link;
	}
	ref.title = r.topic;
	ref.textContent = r.topic;
	newDiv.appendChild(ref);
	document.getElementById("container").appendChild(newDiv);
	ctx.fillStyle=r.fill;
	ctx.fillRect(r.x,r.y,r.width,r.height);
	ctx.strokeStyle="black";
}

function out(s) {
	console.log("2 PAUL THIS IS THE SHAPE:", s);
	if (!ws) {
		return;
	}
	console.log(s);
	ws.send(JSON.stringify(s));
}

// draw a single rect
function circle(c) {
	ctx.fillStyle="lightgreen";
	ctx.beginPath();
	ctx.arc(c.x,c.y,radius,0,Math.PI*2);
	ctx.closePath();
	ctx.fill();
	ctx.fillStyle = "black"; 
	var font = "bold " + radius+"px serif";
	ctx.font = font;
	ctx.textBaseline = "top";
	ctx.fillText(c.id, c.x-(radius-9)/1 ,c.y-(radius-9)/2);
}

// clear the canvas
function clear() {
	$("div.overlay").remove();
	ctx.clearRect(0, 0, WIDTH, HEIGHT);
}

function solicit() {
	for(const property in shapes){
		ratelimit(out,property, 50);
	}
}

function drawsend() {
	console.log(ws);
	if (!ws) {
		return;
	}

	if(!(ws.readyState === 1)) {
		return;
	}
	for(const property in shapes){
		ratelimit(out,property, 50);
	}
	draw();
}

// redraw the scene
function draw() {
	clear();
	ctx.drawImage(background,-130,-60, 1600, 950);   
	// redraw each shape in the shapes[] array
	for(const property in shapes){
		// decide if the shape is a rect or circle
		// (it's a rect if it has a width property)
		if(shapes[property].width){
			rect(shapes[property]);
		}else{
			circle(shapes[property]);
		};

		if (shapes[property].topic) {
			rect({topic: shapes[property].topic, link: shapes[property].link, x: shapes[property].x-(1.5 * radius), y: shapes[property].y-80, width: 120, height: 60, fill: "#b6e3df"})
		}
	}
}

// handle mousedown events
function myDown(e){

	// tell the browser we're handling this mouse event
	e.preventDefault();
	e.stopPropagation();

	// get the current mouse position
	var mx=parseInt(e.clientX-offsetX);
	var my=parseInt(e.clientY-offsetY);

	// test each shape to see if mouse is inside
	dragok=false;
	for(const property in shapes){
		var s=shapes[property];
		// decide if the shape is a rect or circle               
		if(s.width){
			// test if the mouse is inside this rect
			if(mx>s.x && mx<s.x+s.width && my>s.y && my<s.y+s.height){
				// if yes, set that rects isDragging=true
				dragok=true;
				s.isDragging=true;
			}
		}else{
			var dx=s.x-mx;
			var dy=s.y-my;
			// test if the mouse is inside this circle
			if(dx*dx+dy*dy<radius*radius){
				dragok=true;
				s.isDragging=true;
			}
		}
	}

	// save the current mouse position
	startX=mx;
	startY=my;
}

// handle mouseup events
function myUp(e){
	// tell the browser we're handling this mouse event
	e.preventDefault();
	e.stopPropagation();

	// clear all the dragging flags
	dragok = false;
	for(const property in shapes){
		shapes[property].isDragging=false;
	}
}


// handle mouse moves
function myMove(e){
	// if we're dragging anything...
	if (dragok){

		// tell the browser we're handling this mouse event
		e.preventDefault();
		e.stopPropagation();

		// get the current mouse position
		var mx=parseInt(e.clientX-offsetX);
		var my=parseInt(e.clientY-offsetY);

		// calculate the distance the mouse has moved
		// since the last mousemove
		var dx=mx-startX;
		var dy=my-startY;

		// move each rect that isDragging 
		// by the distance the mouse has moved
		// since the last mousemove
		for(const property in shapes){
			var s=shapes[property];
			if(s.isDragging){
				s.x+=dx;
				s.y+=dy;
			}
		}

		// redraw the scene with the new rect positions
		drawsend();

		// reset the starting mouse position for the next mousemove
		startX=mx;
		startY=my;
	}
}


function join(name) {
	shapes[name] = {id: name, x:200, y:100, isDragging: false};
	sleep(500).then(() => {
		drawsend();
	});
}

function ratelimit(f, i, time) {
	const now = +new Date();
	if (!(i in last) || (now - last[i] > time)) {
		last[i] = now;
		f(shapes[i]);
	}
};


function sleep (time) {
	return new Promise((resolve) => setTimeout(resolve, time));
}


function open() {
	if (ws) {
		return false;
	}
	ws = new WebSocket("ws:\/\/" + location.hostname+":"+location.port+"\/echo");
	ws.onopen = function(evt) {
		printer("OPEN");
	}
	ws.onclose = function(evt) {
		printer("CLOSE");
		ws = null;
	}
	ws.onmessage = function(evt) {
		var data = JSON.parse(evt.data);
		if (data.message == "solicit"){
			solicit();
			return;
		}
		if (!(data.id in shapes)) {
			shapes[data.id] = data;
		}
		if (shapes[data.id].isDragging){
			return
		}
		var s = shapes[data.id];
		s.x = data.x;
		s.y = data.y;
		draw();

	}
	ws.onerror = function(evt) {
		printer("ERROR: " + evt.data);
	}
};
