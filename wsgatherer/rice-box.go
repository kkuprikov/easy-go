package wsgatherer

import (
	"time"

	"github.com/GeertJohan/go.rice/embedded"
)

func init() {

	// define files
	file2 := &embedded.EmbeddedFile{
		Filename:    "index.html",
		FileModTime: time.Unix(1571991607, 0),

		Content: string("<html>\n  <head></head>\n  <body>\n    <script type=\"text/javascript\">\n      var sock = null;\n      var wsuri = \"ws://127.0.0.1:1234/ws/subscribe/spectators/123\";\n\n      window.onload = function() {\n        console.log(\"onload\");\n\n        sock = new WebSocket(wsuri);\n\n        sock.onopen = function() {\n          console.log(\"connected to \" + wsuri);\n        };\n\n        sock.onclose = function(e) {\n          console.log(\"connection closed (\" + e.code + \")\");\n        };\n\n        sock.onmessage = function(e) {\n          console.log(\"message received: \" + e.data);\n        };\n      };\n\n      function send() {\n        var msg = document.getElementById(\"message\").value;\n        sock.send(msg);\n      }\n    </script>\n    <h1>WebSocket Echo Test</h1>\n    <form>\n      <p>\n        Message:\n        <input\n          id=\"message\"\n          type=\"text\"\n          value='{\"id\": \"123\", \"event\": \"heatmap\"}'\n        />\n      </p>\n    </form>\n    <button onclick=\"send();\">Send Message</button>\n  </body>\n</html>\n"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1572537627, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file2, // "index.html"

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`static`, &embedded.EmbeddedBox{
		Name: `static`,
		Time: time.Unix(1572537627, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": dir1,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"index.html": file2,
		},
	})
}
