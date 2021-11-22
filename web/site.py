from flask import Flask, render_template, redirect
import json

app = Flask(__name__)

service = """
    <div class="col m-4 serviceCard p-3">
	<div class="row text-center">
		<p class="fw-bold">##NAME##</p>
	</div>
	<div class="row text-center">
		<div class="col">
			<p class="protocolInfo">##PROTOCOL##</p>
		</div>
		<div class="col">
			<p>##MODE##</p>
		</div>
		<div class="col">
			<p class="portInfo">##PORT##</p>
		</div>
	</div>
	<ul>
		##RULE##
	</ul>
	<div class="row">
		<form action="">
			<input type="text" name="" id="s1" placeholder="New rule">
		</form>
	</div>
</div>
"""

@app.route('/')
def hello():
    with open("../config.json", "r") as f:
        data = json.load(f)
    services = ""
    for k in data['services']:
        if k['mode'] == "b":
            mode = "BLACKLIST"
        elif k['mode'] == "w":
            mode = "WHITELIST"
        services += service.replace("##NAME##", k['name']).replace("##PROTOCOL##", k['protocol'].upper()).replace("##PORT##", str(k['dport'])).replace("##MODE##", mode)
        rule = ""
        for r in k['regexList']:
            rule += f"<li>{r[1:-1]}</li>"
        services = services.replace("##RULE##", rule)

    return render_template('index.html', services=services)
    


if __name__ == "__main__":
    app.run(debug=False, host='0.0.0.0', port=8080)