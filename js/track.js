var start = new Date().getTime();
var original_onload = null;
var original_onerror = null;
var original_onbeforeunload = null;
var tracker = null;
var host = null;
var cookieKey = "gbtracker"; // set to null to disable returning user tracking
var thewindow = window;

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

// Modified from: http://stackoverflow.com/questions/105034/create-guid-uuid-in-javascript
generateUUID = function() {
    var d = start;
    var uuid = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
        var r = (d + Math.random()*16)%16 | 0;
        d = Math.floor(d/16);
        return (c=='x' ? r : (r&0x3|0x8)).toString(16);
    });
    return uuid;
}

cookie_tracking = function() {
    if (docCookies.getItem(cookieKey)) {
        return false;
    }
    docCookies.setItem(cookieKey, generateUUID(), Infinity);
    return true;
}

send_data = function(type, data) {
    data["ref"] = thewindow.document.referrer;
    data["lang"] = window.navigator.userLanguage || window.navigator.language;
    data["screen_width"] = screen.width;
    data["screen_height"] = screen.height;
    data["browser_time"] = new Date();
    if (cookieKey != null) {
        data["new_user"] = cookie_tracking();
        data["user_id"] = docCookies.getItem(cookieKey);
    }
    data["uri"] = thewindow.location.pathname;
    data["host"] = thewindow.location.href;
    pars = "";
    for (k in data){
        pars = pars + "&" + k + "=" + encodeURIComponent(data[k])
    }
    // I don't understand the role of tracker variable :(
    // (new Image()).src = host + tracker + '/?t=' + encodeURIComponent(type)+ pars;
    (new Image()).src = host '?t=' + encodeURIComponent(type)+ pars;
}

collect_onbeforeunload_time = function() {
    now = new Date().getTime();
    d = new Array();
    d["beforeunloadtime"] = now;
    d["start"] = start;
    d["total_elapsed_time"] = now - start;
    try {
        if (thewindow.performance != null) {
            n = thewindow.performance.navigation;
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
            send_data("u", d);
        }
    } catch(e) {
        d = {"msg": "tracker exception: " + e.message, "url": thewindow.location.href, "line": e.lineNumber, "file": e.fileName};
        send_data("e", d);
    }
    if (original_onbeforeunload != null) original_onbeforeunload();
}

collect_performance_data = function() {
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
    } catch(e) {
    	  d = {"msg": "tracker exception: " + e.message, "url": thewindow.location.href, "line": e.lineNumber, "file": e.fileName};
    	  send_data("e", d);
    }
    if (original_onload != null) original_onload();
}

collect_errors = function(msg, url, line) {
    d = {"msg":msg, "url": url, "line": line};
    send_data("e", d);
    if (original_onerror != null) original_onerror(msg, url, line);
}

lazy_collect = function() {
    setTimeout(function() {
        collect_performance_data();
    }, 0);
}

gather_geo_info = function() {
    try {
        if (navigator.geolocation) {
            navigator.geolocation.getCurrentPosition(
                function(position) {
                    send_data("g", {"lat": position.coords.latitude, "long": position.coords.longitude});
                }, function(error) {});
        }
    } catch(e) {
    	  d = {"msg": "tracker exception: " + e.message, "url": thewindow.location.href, "line": e.lineNumber, "file": e.fileName};
    	  send_data("e", d);
    }
}

activate = function(h, id, sandbox) {
    if (sandbox) {
        thewindow = window.top;
    } else {
        thewindow = window;
    }
    host = h;
    tracker = id;
    geo = false;

    if (thewindow.onerror != null) original_onerror = thewindow.onerror;
    if (thewindow.onload != null) original_onload = thewindow.onload;
    if (thewindow.onbeforeunload != null) original_onbeforeunload = thewindow.onbeforeunload

    thewindow.onload = lazy_collect
    thewindow.onerror = collect_errors;
    thewindow.onbeforeunload = collect_onbeforeunload_time;
    if (geo == true) { gather_geo_info(); }
}
