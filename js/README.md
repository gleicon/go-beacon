# Simple tracker based on Navigation Timing API

A simple js tracker that collects Navigation Timing API, timers and js exceptions. Events are intercepted to collect data and all calculations are made client-side.

It inserts a 1x1px image into the page and have to be loaded at the `<head>` section or inside a iframe sandbox.

All javascript events are stubbed (i.e. if you already have a function mapped to window.onload it will be called after the tracker

## Why not boomerang.js

Boomerang is pretty good and you should probably be using it. 

This code was made to work with modern browsers and collect specific metrics without plugins or build stage.

It's also smaller than a regular boomerang build.

## Query string and parameters

This tracker returns a couple of events by query string (t= parameter)
    
    u - on before unload metrics: beforeunloadtime, start, total_elapsed_time, type_navigate(RELOAD, BACK_FORWARD, NAVIGATE)
    p - performance data: page_load_time, tcp_time, dns_time, processing_time, processing_time
    e - errors: JSON containing: {"msg":error msg, "url": url, "line": line that the error ocurred}
    g - geolocation info (if geo == true on activate()): JSON object containing 
            {"lat": position.coords.latitude, "long": position.coords.longitude} (it will ask for user permission)

All requests bundle the followind data:
    ref, lang, screen_width, screen_height, browser_time, new_user, url, host

    "ref" = document.referrer;
    "lang" = window.navigator.userLanguage || window.navigator.language;
    "screen_width" = screen.width;
    "screen_height" = screen.height;
    "browser_time" = new Date();
    // disable cookie tracking by commenting the line below
    "new_user" = cookie_tracking();
    "uri" = window.location.pathname;
    "host" = encodeURIComponent(window.location.href);


## Cookies

Cookies are used for marking returning user. Code is based on mozilla code (https://developer.mozilla.org/en-US/docs/Web/API/document/cookie).

You can disable that by setting to `null` the variable `cookieKey` track.js.

## Add to your page, at the `<head>` section

```html
<script type="text/javascript" src="http://location.of.track.js/js/track.js"></script>
<script>
    activate("http://location.of.your.beacon/t.gif", "tracker_name");
</script>
```

## Tracking from inside a iframe sandbox

You can choose to add the script inside a iframe sandbox to prevent any blocking and javascript errors from affecting your page. For that you may add the following iframe snippet inside the pages your're tracking:

```html
<iframe sandbox="allow-same-origin allow-scripts allow-popups allow-forms"
        src="http://location.of.track.html/track.html"
        style="border: 0; width:0px; height:0px;"></iframe>
```

This iframe inserts a `track.html` that need to be something like:

```html
<html>
  <head>
    <script type="text/javascript" src="http://location.of.track.js/js/track.js"></script>
    <script>
        activate("http://location.of.your.beacon/t.gif", "tracker_name", true);
    </script>
  </head>
  <body></body>
</html>
```
