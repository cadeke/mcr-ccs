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
  const [data, setData] = useState<User>();
  const URL_PRD = "http://10.8.8.1:8080";

  // Functions
  const getData = () => {
    const url = URL_PRD;

    // setData(undefined);
    fetch(`${url}/user?q=${input}`)
      .then((response) => response.json())
      .then((json) => setData(json))
      .catch((error) => console.log("ERROR", error));
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
        {data ? (
          <div>
            <table>
              <thead>
                <th>ID</th>
                <th>Username</th>
                <th>Email</th>
              </thead>
              <tbody>
                <tr>
                  <td>{data?.id}</td>
                  <td>{data?.name}</td>
                  <td>{data?.email}</td>
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
