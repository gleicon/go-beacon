var start = new Date().getTime();
var original_onload = null;
var original_onerror = null;
var original_onbeforeunload = null;
var tracker = null;
var host = null;

// From https://developer.mozilla.org/en-US/docs/Web/API/document.cookie
var docCookies = {
  getItem: function (sKey) {
    if (!sKey) { return null; }
    return decodeURIComponent(document.cookie.replace(new RegExp("(?:(?:^|.*;)\\s*" + encodeURIComponent(sKey).replace(/[\-\.\+\*]/g, "\\$&amp;") + "\\s*\\=\\s*([^;]*).*$)|^.*$"), "$1")) || null;
  },
  setItem: function (sKey, sValue, vEnd, sPath, sDomain, bSecure) {
    if (!sKey || /^(?:expires|max\-age|path|domain|secure)$/i.test(sKey)) { return false; }
    var sExpires = "";
    if (vEnd) {
      switch (vEnd.constructor) {
        case Number:
          sExpires = vEnd === Infinity ? "; expires=Fri, 31 Dec 9999 23:59:59 GMT" : "; max-age=" + vEnd;
          break;
        case String:
          sExpires = "; expires=" + vEnd;
          break;
        case Date:
          sExpires = "; expires=" + vEnd.toUTCString();
          break;
      }
    }
    document.cookie = encodeURIComponent(sKey) + "=" + encodeURIComponent(sValue) + sExpires + (sDomain ? "; domain=" + sDomain : "") + (sPath ? "; path=" + sPath : "") + (bSecure ? "; secure" : "");
    return true;
  },
};

if (window.onerror != null) original_onerror = window.onerror;
if (window.onload != null) original_onload = window.onload;
if (window.onbeforeunload != null) original_onbeforeunload = window.onbeforeunload

cookie_tracking = function() {
    if (docCookies.getItem("gbtracker")) {
        return false;
    }
    docCookies.setItem("gbtracker", 1, Infinity);
}
 

send_data = function(type, data) {
    data["ref"] = document.referrer;
    data["lang"] = window.navigator.userLanguage || window.navigator.language;
    data["screen_width"] = screen.width;
    data["screen_height"] = screen.height;
    data["browser_time"] = new Date();
    // disable cookie tracking by commenting the line below
    data["new_user"] = cookie_tracking();
    data["uri"] = window.location.pathname; 
    data["host"] = encodeURIComponent(window.location.href);
    pars = "";
    for (k in data){
        pars = pars + "&"+k+"="+e(data[k])
    }
    (new Image()).src = host + tracker + '/?t=' + e(type)+ pars;
}

collect_onbeforeunload_time = function(){
            now = new Date().getTime();
            d = new Array();
            d["beforeunloadtime"] = now;
            d["start"] = start;
            d["total_elapsed_time"] = now - start;
            try {
                if (performance != null){
                    n = performance.navigation;
                    switch(n.type){
                        case n.TYPE_RELOAD:
                            d["type_navigate"] = "reload";
                            break;
                        case n.TYPE_BACK_FORWARD:
                            d["type_navigate"] = "back_forward";
                            break;
                        default:
                        case n.TYPE_NAVIGATE:
                            d["type_navigate"] = "navigate";
                            break;
                    }
                    send_data("u", d)
                }
            } catch(e){
		d = {"msg": "tracker exception: " + e.message, "url": window.location.href, "line": e.lineNumber, "file": e.fileName};
		send_data("e", d);
            }
        if (original_onbeforeunload != null) original_onbeforeunload();
}

collect_performance_data = function(){
    var d;
    try {
        if (performance != null) {
            var t = performance.timing;
            var n = performance.navigation;
            if (t.loadEventEnd > 0) {
              var page_load_time = t.loadEventEnd - t.navigationStart;
              var tcp_time = t.connectEnd - t.connectStart;
              var dns_time = t.domainLookupEnd - t.domainLookupStart; 
              var processing_time = t.loadEventEnd - t.domLoading
              var d = new Array();
              if (n.type == n.TYPE_NAVIGATE || n.type == n.TYPE_RELOAD) {
                d["page_load_time"] = page_load_time;
                d["tcp_time"] = tcp_time;
                d["dns_time"] = dns_time;
                d["processing_time"] = processing_time;
              }
            }
        } else {
            now = new Date().getTime();
            var processing_time = now - start;
            d["processing_time"] = processing_time;
        }
        send_data("p", d);
    } catch(e){
    	d = {"msg": "tracker exception: " + e.message, "url": window.location.href, "line": e.lineNumber, "file": e.fileName};
    	send_data("e", d);
    }
    if (original_onload != null) original_onload();
}

collect_errors = function(msg, url, line) {
    d = {"msg":msg, "url": url, "line": line};
    send_data("e", d);
    if (original_onerror != null) original_onerror(msg, url, line);
}

lazy_collect = function(){ setTimeout(function(){
        collect_performance_data();
    }, 0);
}

gather_geo_info = function(){
    try {
        if (navigator.geolocation){
            navigator.geolocation.getCurrentPosition(
                    function(position){ 
                        send_data("g", {"lat": position.coords.latitude, "long": position.coords.longitude}); 
                    }, function(error){});
        }
    } catch(e){
    	d = {"msg": "tracker exception: " + e.message, "url": window.location.href, "line": e.lineNumber, "file": e.fileName};
    	send_data("e", d);
    }
}

activate = function(h, id) {
    host = h;
    tracker = id;
    geo = false;
    window.onload = lazy_collect
    window.onerror = collect_errors;
    window.onbeforeunload = collect_onbeforeunload_time;
    if (geo == true) { gather_geo_info(); }
}
