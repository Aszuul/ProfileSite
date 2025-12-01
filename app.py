from flask import Flask, render_template, url_for
import requests
import os

CAPTCHA_SECRET_KEY=os.environ['CAPTCHA_SECRET_KEY']
CAPTCHA_SITE_KEY=os.environ['CAPTCHA_SITE_KEY']

app = Flask(__name__)
app.config['SECRET_KEY'] = os.environ.get('SECRET_KEY', 'very secret')
FLASK_DEBUG = os.environ.get('FLASK_ENV', True)

def verify_recaptcha(token, action):
    """Verifies the reCAPTCHA token with Google."""
    RECAPTCHA_VERIFY_URL = 'https://www.google.com/recaptcha/api/siteverify'
    
    response = requests.post(RECAPTCHA_VERIFY_URL, data={
        'secret': CAPTCHA_SECRET_KEY,
        'response': token,
        'remoteip': requests.request.remote_addr # Optional but recommended
    })
    
    result = response.json()
    
    # Check if verification succeeded and the action name matches
    if not result.get('success'):
        return False, "Verification failed (Google response failed)."

    # Check the score and the action name
    MIN_SCORE = 0.5 # Set your acceptable threshold (0.5 is default)
    if result.get('score', 0.0) < MIN_SCORE:
        return False, f"Score too low: {result.get('score')} < {MIN_SCORE}"
        
    if result.get('action') != action:
        # This prevents an attacker from using a token from a different action (e.g., 'login')
        return False, f"Action mismatch: expected '{action}', got '{result.get('action')}'"
        
    return True, "Success"


@app.route('/')
def home():
    return render_template('index.html')

@app.route('/about')
def about():
    return render_template('about.html')

@app.route('/projects')
def projects():
    return render_template('projects.html')

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=80, debug=FLASK_DEBUG)