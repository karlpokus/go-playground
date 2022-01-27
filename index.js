document.addEventListener("DOMContentLoaded", function() {
  console.log("dom loaded. Adding event listeners")
  document.getElementById('submit').onclick = function() {
    console.log("user clicked button")
    v = document.getElementById("code").value;
    httpRequest = new XMLHttpRequest();
    httpRequest.onreadystatechange = function() {
      try {
        if (httpRequest.readyState === XMLHttpRequest.DONE) {
          if (httpRequest.status === 200) {
            document.getElementById('result').innerHTML = httpRequest.responseText;
          } else {
            console.error('request failed');
          }
        }
      }
      catch( e ) {
        console.error('Caught Exception: ' + e.description);
      }
    };
    httpRequest.open('POST', '/code');
    httpRequest.send(v);
  }
});
