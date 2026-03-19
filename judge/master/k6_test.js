import http from 'k6/http';
import { check } from 'k6';

export const options = {
  vus: 10,
  iterations: 200,
};

const payloads = [
  {
    QuestionId: 1069,
    runtime: 'c++17',
    code: `#include <bits/stdc++.h>\r
#include<string>\r
using namespace std;\r
#define int long long\r
 \r
signed main()\r
{ios_base::sync_with_stdio(false);\r
cin.tie(NULL);cout.tie(NULL);\r
 \r
 {\r
 \r
    string a;\r
    cin>>a;\r
    int maxi=1;\r
    int pt=0;\r
    int n= a.length();\r
 \r
    for(int i=1;i<n;i++){\r
        if(a[i]!=a[i-1]){\r
            \r
            maxi=max(maxi,i-pt);\r
            pt=i;\r
        }\r
 \r
    }\r
    maxi=max(maxi,n-pt);\r
 \r
    cout<<maxi<<endl;\r
 \r
 }\r
 \r
return 0;\r
}`,
  },
  {
    QuestionId: 1069,
    runtime: 'python3',
    code: `import sys\r
\r
s = sys.stdin.readline().strip()\r
\r
max_len = 1\r
curr_len = 1\r
\r
for i in range(1, len(s)):\r
    if s[i] == s[i-1]:\r
        curr_len += 1\r
    else:\r
        curr_len = 1\r
    max_len = max(max_len, curr_len)\r
\r
print(max_len)`,
  },
  {
    QuestionId: 1069,
    runtime: 'node-25',
    code: `const fs = require("fs");\r
\r
const s = fs.readFileSync(0, "utf8").trim();\r
\r
let maxLen = 1;\r
let curr = 1;\r
\r
for (let i = 1; i < s.length; i++) {\r
    if (s[i] === s[i - 1]) {\r
        curr++;\r
    } else {\r
        curr = 1;\r
    }\r
    if (curr > maxLen) maxLen = curr;\r
}\r
\r
console.log(maxLen);`,
  },
  {
    QuestionId: 1069,
    runtime: 'c++23',
    code: `#include <bits/stdc++.h>\r
#include<string>\r
using namespace std;\r
#define int long long\r
 \r
signed main()\r
{ios_base::sync_with_stdio(false);\r
cin.tie(NULL);cout.tie(NULL);\r
 \r
 {\r
 \r
    string a;\r
    cin>>a;\r
    int maxi=1;\r
    int pt=0;\r
    int n= a.length();\r
 \r
    for(int i=1;i<n;i++){\r
        if(a[i]!=a[i-1]){\r
            \r
            maxi=max(maxi,i-pt);\r
            pt=i;\r
        }\r
 \r
    }\r
    maxi=max(maxi,n-pt);\r
 \r
    cout<<maxi<<endl;\r
 \r
 }\r
 \r
return 0;\r
}`,
  },
];

const params = {
  headers: {
    'Content-Type': 'application/json',
  },
};

export default function () {
  const payload = payloads[__ITER % payloads.length];
  const res = http.post(
    'http://localhost:7000/judge/',
    JSON.stringify(payload),
    params
  );

  check(res, {
    'status is 200': (r) => r.status === 202,
  });

  console.log(
    `VU=${__VU} ITER=${__ITER} runtime=${payload.runtime} status=${res.status}`
  );
}
