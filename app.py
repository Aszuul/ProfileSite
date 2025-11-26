from flask import Flask, render_template, url_for
import os

app = Flask(__name__)
app.config['SECRET_KEY'] = os.environ.get('SECRET_KEY', 'very secret')
FLASK_DEBUG = os.environ.get('FLASK_ENV', True)

@app.route('/')
def home():
    return render_template('index.html')

@app.route('/about')
def about():
    return render_template('about.html')

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=80, debug=FLASK_DEBUG)