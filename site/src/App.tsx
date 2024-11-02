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

  // Functions
  const getData = async () => {
    setApiData(undefined);
    const url = "/api/user";
    try {
      console.log("Start api call");
      const res = await fetch(`${url}?q=${input}`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });
      const resData = await res.json();
      console.log("Got data", resData);
      setApiData(resData);
    } catch (err) {
      console.log("Got error", err);
    } finally {
      console.log("End api call");
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
        {true ? (
          <div className="Table">
            <table>
              <thead>
                <th>ID</th>
                <th>Name</th>
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
