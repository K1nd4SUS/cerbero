from flask import Flask, render_template, request, redirect
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
		<form action="/jsonUpdater/##NAME##" method="POST">
			<input type="text" name="regex" placeholder="New rule">
            <input type="submit" value="update">
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
        services += service.replace("##NAME##", k['name'].replace(" ", "")).replace("##PROTOCOL##", k['protocol'].upper()).replace("##PORT##", str(k['dport'])).replace("##MODE##", mode)
        rule = ""
        for r in k['regexList']:
            rule += f"<li>{r[1:-1]}</li>"
        services = services.replace("##RULE##", rule)

    return render_template('index.html', services=services)
    
@app.route('/jsonUpdater/<path:service>', methods=['POST'])
def jsonUpdate(service):
    regex = "(" + request.form.get('regex') + ")"
    with open("../config.json", "r") as f:
        data = json.load(f)
    
    # find the right service
    for k in range(len(data['services'])):
        if data['services'][k]['name'].replace(" ", "") == service:
            index = k

    # check if the rule already exist, in this case it need to be removed
    # otherwise we can add the new rule to the file
    for rule in data['services'][index]['regexList']:
        if(rule == regex):
            #remove
            data['services'][index]['regexList'].remove(rule)
            with open('../config.json', 'w') as json_file:
                json.dump(data, json_file)
                return redirect("/")

    data['services'][index]['regexList'].append(regex)
    with open('../config.json', 'w') as json_file:
        json.dump(data, json_file)
    return redirect("/")



if __name__ == "__main__":
    app.run(debug=True, host='0.0.0.0', port=8080)