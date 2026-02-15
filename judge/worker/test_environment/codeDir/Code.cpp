// #include <bits/stdc++.h>
// #include <iostream>
// #include <chrono>
// #include <thread>
// #include <vector>
// using namespace std; // Using namespace std makes code less verbose
//
// int main() {
//     vector<int> v = {5, 3, 1, 4, 9, 7};
//     // cout << "Original vector: ";
//     // for (int i : v) {
//     //     cout << i << " ";
//     // }
//     // cout << endl;
//     //
//     // sort(v.begin(), v.end()); // sort function included via bits/stdc++.h
//     //
//     // cout << "Sorted vector: ";
//     // for (int i : v) {
//     //     cout << i << " ";
//     // }
//     // cout << endl;
//
//   int x,y;
//
//   std::this_thread::sleep_for(std::chrono::milliseconds(500));
//
//   float t=1.3,k=1.67;
//   float final;
//   vector<long long> vlong;
//   for(int i=0;i<1000000;i++){
//
//     vlong.push_back(i);
//
//   }
//   cin>>x>>y;
//   cout<<x<<" "<<y;
//
//     return 0;
// }
//
//
//

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
