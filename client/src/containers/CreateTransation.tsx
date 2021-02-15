/* eslint-disable camelcase */
/* eslint-disable react/no-unescaped-entities */
/* eslint-disable jsx-a11y/label-has-associated-control */
/* eslint-disable react/jsx-curly-newline */
// eslint-disable-next-line no-use-before-define
import React, { useReducer, useEffect, ChangeEvent } from "react";
import MaskedInput from "react-text-mask";
import createNumberMask from "text-mask-addons/dist/createNumberMask";

/**
 * Application State
 */

type Tranaction = {
  id: string;
  customer_id: string;
  load_amount: string;
  time: string;
};

type State = {
  transaction: Tranaction;
  success: boolean;
  submitted: boolean;
  error: string;
};

const initialState = {
  transaction: {
    id: Math.floor(Math.random() * Math.floor(9999999)).toString(),
    customer_id: "",
    load_amount: "",
    time: new Date().toISOString(),
  },
  success: false,
  submitted: false,
  error: "",
};

type Action =
  | { type: "setTransaction"; payload: State["transaction"] }
  | { type: "setSuccess"; payload: State["success"] }
  | { type: "setSubmitted"; payload: State["submitted"] }
  | { type: "setError"; payload: State["error"] };

function reducer(state: State, action: Action): State {
  switch (action.type) {
    case "setTransaction":
      return {
        ...state,
        transaction: action.payload,
        success: false,
        error: "",
      };
    case "setSuccess":
      return {
        ...state,
        success: action.payload,
      };
    case "setSubmitted":
      return {
        ...state,
        success: false,
        submitted: action.payload,
      };
    case "setError":
      return {
        ...state,
        error: action.payload,
      };
    default:
      throw new Error();
  }
}

/**
 * Create Transation
 */
function CreateTransation(): React.ReactElement {
  const [state, dispatch] = useReducer(reducer, initialState);

  const currencyOptions = {
    prefix: "$",
    suffix: "",
    includeThousandsSeparator: false,
    thousandsSeparatorSymbol: ",",
    allowDecimal: true,
    decimalSymbol: ".",
    decimalLimit: 2,
    integerLimit: 12,
    allowNegative: false,
    allowLeadingZeroes: false,
  };

  const currencyMask = createNumberMask({
    ...currencyOptions,
  });

  const customerOptions = {
    prefix: "",
    suffix: "",
    includeThousandsSeparator: false,
    thousandsSeparatorSymbol: ",",
    allowDecimal: false,
    decimalSymbol: ".",
    decimalLimit: 2,
    integerLimit: 12,
    allowNegative: false,
    allowLeadingZeroes: false,
  };

  const inputMask = createNumberMask({
    ...customerOptions,
  });

  // Form submission
  const processData = (event: React.MouseEvent<HTMLButtonElement>) => {
    event.preventDefault();
    if (state.transaction.customer_id && state.transaction.load_amount) {
      fetch("http://localhost:8000/transaction", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          ...state.transaction,
        }),
      })
        .then((response) => response.json())
        .then((res) => {
          // eslint-disable-next-line no-console
          if (res.accepted) {
            dispatch({ type: "setSubmitted", payload: true });
          } else {
            dispatch({
              type: "setError",
              payload: "Unable to process your transaction.",
            });
          }
        })
        .catch((err) => {
          dispatch({
            type: "setError",
            payload: "Error retrieving interest rate!",
          });
        });
    } else {
      dispatch({
        type: "setError",
        payload: "Form Is Incomplete!",
      });
    }
  };

  return (
    <div tabIndex={-1} role="group">
      <div>
        <div className="css-e63m0r">
          <div className="css-6g7461" />
        </div>
        <div className="css-522n1y">
          <div className="css-18opf0w">
            <div className="css-1hkiohq">
              <img
                alt="KOHO Logo"
                src="https://web.koho.ca/static/media/logo.4a1b90d4.svg"
                className="css-1fld7e1"
              />
            </div>
            {!state.submitted ? (
              <>
                <h1 className="css-1uzylmb">Let’s Add Your Transaction</h1>
                <div className="css-1oxkbvq">
                  Do keep in mind, that there is daily limit of 3 transactions
                  with a max limit of $5000.00 per day or $20,000.00 per week.
                </div>

                {state.error && <h2>{state.error}</h2>}

                <form noValidate autoComplete="off">
                  <div className="css-1aoknzf">
                    <div className="css-79elbk">
                      <div className="MuiFormControl-root MuiTextField-root css-1qj5j9a MuiFormControl-fullWidth">
                        <label
                          className="MuiFormLabel-root MuiInputLabel-root MuiInputLabel-formControl MuiInputLabel-animated MuiInputLabel-outlined Mui-required Mui-required"
                          data-shrink="false"
                        >
                          ID
                          <span className="MuiFormLabel-asterisk MuiInputLabel-asterisk">
                            *
                          </span>
                        </label>
                        <div className="MuiInputBase-root MuiOutlinedInput-root MuiInputBase-fullWidth MuiInputBase-formControl">
                          <input
                            aria-invalid="false"
                            autoComplete="off"
                            name="id"
                            placeholder="Transaction ID"
                            required
                            type="text"
                            className="MuiInputBase-input MuiOutlinedInput-input"
                            value={state.transaction.id}
                            disabled
                          />
                          <fieldset
                            aria-hidden="true"
                            className="jss188 MuiOutlinedInput-notchedOutline"
                          >
                            <legend className="jss189" />
                          </fieldset>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div className="css-1aoknzf">
                    <div className="css-79elbk">
                      <div className="MuiFormControl-root MuiTextField-root css-1qj5j9a MuiFormControl-fullWidth">
                        <label
                          className="MuiFormLabel-root MuiInputLabel-root MuiInputLabel-formControl MuiInputLabel-animated MuiInputLabel-outlined Mui-required Mui-required"
                          data-shrink="false"
                        >
                          Customer ID
                          <span className="MuiFormLabel-asterisk MuiInputLabel-asterisk">
                            *
                          </span>
                        </label>
                        <div className="MuiInputBase-root MuiOutlinedInput-root MuiInputBase-fullWidth MuiInputBase-formControl">
                          <MaskedInput
                            mask={inputMask}
                            className="MuiInputBase-input MuiOutlinedInput-input"
                            placeholder="Customer ID"
                            guide={false}
                            data-testid="customer_id"
                            id="customer_id"
                            onChange={(event) =>
                              dispatch({
                                type: "setTransaction",
                                payload: {
                                  ...state.transaction,
                                  customer_id: event.target.value,
                                },
                              })
                            }
                            value={state.transaction.customer_id}
                            required
                          />
                          <fieldset
                            aria-hidden="true"
                            className="jss188 MuiOutlinedInput-notchedOutline"
                          >
                            <legend className="jss189" />
                          </fieldset>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div className="css-1aoknzf">
                    <div className="css-79elbk">
                      <div className="MuiFormControl-root MuiTextField-root css-1qj5j9a MuiFormControl-fullWidth">
                        <label
                          className="MuiFormLabel-root MuiInputLabel-root MuiInputLabel-formControl MuiInputLabel-animated MuiInputLabel-outlined Mui-required Mui-required"
                          data-shrink="false"
                        >
                          Amount
                          <span className="MuiFormLabel-asterisk MuiInputLabel-asterisk">
                            *
                          </span>
                        </label>
                        <div className="MuiInputBase-root MuiOutlinedInput-root MuiInputBase-fullWidth MuiInputBase-formControl">
                          <MaskedInput
                            mask={currencyMask}
                            className="MuiInputBase-input MuiOutlinedInput-input"
                            placeholder="Amount to load into account"
                            guide={false}
                            data-testid="load_amount"
                            id="load_amount"
                            onChange={(event) =>
                              dispatch({
                                type: "setTransaction",
                                payload: {
                                  ...state.transaction,
                                  load_amount: event.target.value,
                                },
                              })
                            }
                            value={state.transaction.load_amount}
                            required
                          />
                          <fieldset
                            aria-hidden="true"
                            className="jss188 MuiOutlinedInput-notchedOutline"
                          >
                            <legend className="jss189" />
                          </fieldset>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div className="css-1aoknzf">
                    <div className="css-79elbk">
                      <div className="MuiFormControl-root MuiTextField-root css-1qj5j9a MuiFormControl-fullWidth">
                        <label
                          className="MuiFormLabel-root MuiInputLabel-root MuiInputLabel-formControl MuiInputLabel-animated MuiInputLabel-outlined Mui-required Mui-required"
                          data-shrink="false"
                        >
                          Timestamp
                          <span className="MuiFormLabel-asterisk MuiInputLabel-asterisk">
                            *
                          </span>
                        </label>
                        <div className="MuiInputBase-root MuiOutlinedInput-root MuiInputBase-fullWidth MuiInputBase-formControl">
                          <input
                            aria-invalid="false"
                            autoComplete="off"
                            name="time"
                            required
                            type="text"
                            className="MuiInputBase-input MuiOutlinedInput-input"
                            value={state.transaction.time}
                            disabled
                          />
                          <fieldset
                            aria-hidden="true"
                            className="jss188 MuiOutlinedInput-notchedOutline"
                          >
                            <legend className="jss189" />
                          </fieldset>
                        </div>
                      </div>
                    </div>
                  </div>
                  <div>
                    <button
                      type="submit"
                      className="css-1j8penb"
                      onClick={(event) => processData(event)}
                    >
                      Process
                    </button>
                  </div>
                </form>
                <div className="css-h5fkc8">
                  <div className="css-3ec9p3">
                    <svg
                      className="MuiSvgIcon-root"
                      focusable="false"
                      viewBox="0 0 24 24"
                      aria-hidden="true"
                      role="presentation"
                    >
                      <path d="M18 8h-1V6c0-2.76-2.24-5-5-5S7 3.24 7 6v2H6c-1.1 0-2 .9-2 2v10c0 1.1.9 2 2 2h12c1.1 0 2-.9 2-2V10c0-1.1-.9-2-2-2zm-6 9c-1.1 0-2-.9-2-2s.9-2 2-2 2 .9 2 2-.9 2-2 2zm3.1-9H8.9V6c0-1.71 1.39-3.1 3.1-3.1 1.71 0 3.1 1.39 3.1 3.1v2z" />
                    </svg>
                    <div className="css-h4a49s">
                      Your information is securely stored.
                      <a
                        href="https://www.koho.ca/legal#PrivacyPolicy"
                        className="css-4ogr89"
                      >
                        See our Privacy Policy Policy
                      </a>
                    </div>
                  </div>
                </div>
              </>
            ) : (
              <h1 className="css-1uzylmb">Your Transactions was processed!</h1>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

export default CreateTransation;
