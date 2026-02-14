#include <bits/stdc++.h>
using namespace std; // Using namespace std makes code less verbose

int main() {
    vector<int> v = {5, 3, 1, 4, 9, 7};
    cout << "Original vector: ";
    for (int i : v) {
        cout << i << " ";
    }
    cout << endl;

    sort(v.begin(), v.end()); // sort function included via bits/stdc++.h

    cout << "Sorted vector: ";
    for (int i : v) {
        cout << i << " ";
    }
    cout << endl;

    return 0;
}

