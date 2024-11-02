import { useState } from "react";
import "./App.css";

function App() {
  // Types
  interface User {
    id: number;
    name: string;
    email: string;
  }

  // States
  const [input, setInput] = useState("");
  const [apiData, setApiData] = useState<User>();
  const URL_PRD = "http://10.8.8.1:8080";

  // Functions
  const getData = async () => {
    const url = URL_PRD;

    // setData(undefined);
    try {
      console.log("Start api call");
      const res = await fetch(`${url}/user?q=${input}`);
      const resData = await res.json();
      console.log("Got data", resData);
      setApiData(resData);
    } catch (err) {
      console.log("Got error", err);
    }
  };

  return (
    <>
      <div className="App">
        <h1>User Search</h1>
        <input
          type="text"
          value={input}
          onChange={(e) => setInput(e.target.value)}
        />
        <button onClick={() => getData()}>Search</button>
        {apiData ? (
          <div>
            <table>
              <thead>
                <th>ID</th>
                <th>Username</th>
                <th>Email</th>
              </thead>
              <tbody>
                <tr>
                  <td>{apiData?.id}</td>
                  <td>{apiData?.name}</td>
                  <td>{apiData?.email}</td>
                </tr>
              </tbody>
            </table>
          </div>
        ) : null}
      </div>
    </>
  );
}

export default App;
