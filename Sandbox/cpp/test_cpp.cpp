#include <iostream>
#include <queue>

using namespace std;
void ShowQueue(queue<int>a)
{
	queue<int>temp = a;
	while (!temp.empty())
	{
		cout << temp.front() << " ";
		temp.pop();
	}
}
int main()
{
	int t;
	cin >> t;
	while (t--)
	{
		int n, k;
		cin >> n >> k;
		queue<int>list;
		for (int i = 1;i <= n;i++)
		{
			list.push(i);
		}
		queue<int>out;
		int i = 1;
		while (out.size() < n - 1)
		{
			if (i == k)
			{
				out.push(list.front());
				list.pop();
				i = 1;
			}
			else
			{
				i++;
				list.push(list.front());
				list.pop();
			}
		}
		ShowQueue(out);
		cout << list.front() << endl;
	}
}