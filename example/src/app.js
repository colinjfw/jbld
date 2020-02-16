import * as foobar from "foobar";

 export class App {
  hello() { return "hello" }
  mod() { return import("./mod") }
}
