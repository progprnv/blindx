import re
import sys
import html
import base64
import urllib.parse
import requests

BANNER = r'''   ___      __    _               __   _  __  
  / _ )    / /   (_)   ___    ___/ /  | |/_/  
 / _  |   / /   / /   / _ \\  / _  /  _>  <   
/____/   /_/   /_/   /_//_/  \,_/  /_/|_|      
                                            
     ----- blindx V1.0 by progprnv | pr0gx | Pranav Jayan '''

# Encoding functions
def html_encode(s, times):
    for _ in range(times):
        s = html.escape(s)
    return s

def url_encode(s, times):
    for _ in range(times):
        s = urllib.parse.quote(s)
    return s

def js_escape(s, times):
    for _ in range(times):
        s = s.replace('\\', '\\\\')
        s = s.replace("'", "\\'")
        s = s.replace('"', '\\"')
    return s

def unicode_escape(s, times):
    for _ in range(times):
        s = ''.join(f"\\u{ord(c):04x}" for c in s)
    return s

def base64_encode(s, times):
    for _ in range(times):
        s = base64.b64encode(s.encode()).decode()
    return s


def parse_raw_request(raw):
    parts = raw.split('\r\n\r\n', 1)
    header_lines = parts[0].split('\r\n')
    body = parts[1] if len(parts) > 1 else ''
    method, path, _ = header_lines[0].split()

    headers = {}
    for line in header_lines[1:]:
        k, v = line.split(':', 1)
        headers[k.strip()] = v.strip()

    host = headers.get('Host')
    url = f"http://{host}{path}" if not path.startswith('http') else path
    return method, url, headers, body


def main():
    print(BANNER)
    print('\nPaste your raw POST (with CRLF) and end with an empty line:')
    lines = []
    while True:
        line = sys.stdin.readline()
        if line.strip() == '':
            break
        lines.append(line.rstrip('\n'))
    raw = '\r\n'.join(lines)

    try:
        method, url, headers, body = parse_raw_request(raw)
    except Exception as e:
        print(f"[!] Failed to parse request: {e}")
        sys.exit(1)

    # Params selection
    params = []
    while True:
        p = input('Enter parameter to inject payload into: ').strip()
        params.append(p)
        if input('Add another parameter? (Y/N): ').lower() != 'y':
            break

    payload = input('Enter your payload: ')

    # Encoding options
    options = [f"{i+1}) {name}" for i, name in enumerate([
        'HTML encode x1','HTML encode x2','HTML encode x3',
        'URL encode x1','URL encode x2','URL encode x3',
        'JS escape x1','JS escape x2','JS escape x3',
        'Unicode escape x1','Unicode escape x2','Unicode escape x3',
        'Base64 encode x1','Base64 encode x2','Base64 encode x3',
        'All variants','No encode'
    ])]
    print('\n'.join(options))
    choice = int(input('Select option number (1-17): '))

    variants = []
    if 1 <= choice <= 15:
        funcs = [html_encode, html_encode, html_encode,
                 url_encode, url_encode, url_encode,
                 js_escape, js_escape, js_escape,
                 unicode_escape, unicode_escape, unicode_escape,
                 base64_encode, base64_encode, base64_encode]
        variants.append(funcs[choice-1](payload, ((choice-1)%3)+1))
    elif choice == 16:
        for i, fn in enumerate([html_encode]*3 + [url_encode]*3 + [js_escape]*3 + [unicode_escape]*3 + [base64_encode]*3):
            variants.append(fn(payload, (i%3)+1))
    else:
        variants.append(payload)

    # Header injection
    add_headers = {}
    if input('Inject headers? (Y/N): ').lower() == 'y':
        while True:
            hn = input('Header name: ').strip()
            hv = input('Header value (use {{payload}}): ').strip()
            add_headers[hn] = hv
            if input('Add another header? (Y/N): ').lower() != 'y':
                break

    session = requests.Session()
    for i, var in enumerate(variants, start=1):
        data = dict(urllib.parse.parse_qsl(body))
        for p in params:
            data[p] = var
        h = headers.copy()
        for hn, hv in add_headers.items():
            h[hn] = hv.replace('{{payload}}', var)
        resp = session.request(method, url, headers=h, data=data)
        print(f"[{i}] {url} -> {resp.status_code}")

if __name__ == '__main__':
    main()
