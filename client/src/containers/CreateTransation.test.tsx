import {
  render,
  screen,
  fireEvent,
  waitFor,
  waitForElementToBeRemoved,
} from "@testing-library/react";
// eslint-disable-next-line import/no-extraneous-dependencies
import fetch from "jest-fetch-mock";
import CreateTransation from "./CreateTransation";

test("Customer ID Field", () => {
  render(<CreateTransation />);
  const input = screen.getByTestId("customer_id") as HTMLInputElement;

  // Expect - field to be empty
  expect(input.value).toBe("");

  fireEvent.change(input, { target: { value: "1000000" } });

  // Expect - field to be 1000000 when typed '1000000'
  expect(input.value).toBe("1000000");

  fireEvent.change(input, { target: { value: "aa9999999a" } });

  // Expect - field to be 9999999 when typed 'aa9999999a'
  expect(input.value).toBe("9999999");

  fireEvent.change(input, { target: { value: "" } });
});

test("Amount Field", () => {
  render(<CreateTransation />);
  const amount = screen.getByTestId("load_amount") as HTMLInputElement;

  // Expect - field to be empty
  expect(amount.value).toBe("");

  fireEvent.change(amount, { target: { value: "1000000" } });

  // Expect - field to be $1000000 when typed '1000000'
  expect(amount.value).toBe("$1000000");

  fireEvent.change(amount, { target: { value: "aaa99.99a" } });

  // Expect - field to be $99.9 when typed 'aaa99.99a'
  expect(amount.value).toBe("$99.99");
});

test("Submit Transation", async () => {
  const mockSuccessResponse = { accepted: true };
  const mockJsonPromise = Promise.resolve(mockSuccessResponse);
  const mockFetchPromise = Promise.resolve({
    json: () => mockJsonPromise,
  });
  const globalRef: any = global;
  globalRef.fetch = jest.fn().mockImplementation(() => mockFetchPromise);

  render(<CreateTransation />);

  const customer = screen.getByTestId("customer_id") as HTMLInputElement;
  const amount = screen.getByTestId("load_amount") as HTMLInputElement;

  fireEvent.click(screen.getByTestId("button-submit"));

  // Expect form error
  await waitFor(() =>
    expect(screen.getByText("Form Is Incomplete!")).toBeDefined()
  );

  fireEvent.change(customer, { target: { value: "999" } });
  fireEvent.click(screen.getByTestId("button-submit"));

  // Expect form error
  await waitFor(() =>
    expect(screen.getByText("Form Is Incomplete!")).toBeDefined()
  );

  fireEvent.change(amount, { target: { value: "1000000" } });
  fireEvent.click(screen.getByTestId("button-submit"));

  // Expect form error
  await waitFor(() =>
    expect(screen.getByText("Your Transactions was processed!")).toBeDefined()
  );
});
