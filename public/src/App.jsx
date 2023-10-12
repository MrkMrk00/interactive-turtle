import { createSignal } from "solid-js";

function App() {
  const [data, setData] = createSignal();

  function getData(ev) {
    ev.preventDefault();
    const data = Object.fromEntries(new FormData(ev.currentTarget));

    const promise = fetch('http://localhost:8069/intoTripples', {
      method: 'POST',
      body: data.turtle,
    });

    promise
      .then(r => r.json())
      .then(JSON.stringify)
      .then(setData)
      .catch(console.error);
  }

  return (
    <div>
      <form 
        onSubmit={getData}
        class="p-4 flex flex-col gap-2 max-w-xl items-start"
      >
        <textarea 
          id="turtle-textarea"
          class="border-2 focus:border-green-500 shadow-md rounded-2xl overflow-hidden p-4 outline-none"
          cols="60"
          rows="10"
          name="turtle"
          placeholder="@prefix : <http://tvoje.mama/> ."
        ></textarea>
        <button 
          type="submit"
          class="rounded-2xl border-2 bg-white border-green-500 p-4 shadow-md font-bold hover:brightness-95"
        >Konvertovat</button>
      </form>

      <code>{data()}</code>
    </div>
  );
}

export default App;
