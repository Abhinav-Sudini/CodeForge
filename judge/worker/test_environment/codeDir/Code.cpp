#include <bits/stdc++.h>
#include <string>
using namespace std;
#define int long long

signed main() {
  ios_base::sync_with_stdio(false);
  cin.tie(NULL);
  cout.tie(NULL);

  {

    string a;
    cin >> a;

    int n = a.length();
    map<char, int> m;
    for (int i = 0; i < n; i++) {
      m[a[i]]++;
    }
    int flag = 0;
    char temp;
    string out;
    for (auto &x : m) {
      for (int i = 0; i < x.second / 2; i++)
        out += x.first;

      if (x.second & 1) {
        temp = x.first;
        flag++;
      }
    }
    if (flag == 0) {
      cout << out;
      reverse(out.begin(), out.end());
      cout << out << endl;
    } else if (flag == 1) {
      cout << out << temp;
      reverse(out.begin(), out.end());
      cout << out << endl;
    } else {
      cout << "NO SOLUTION\n";
    }
  }

  return 0;
}