import React from 'react';
import { Fieldset, Field, createValue, Input } from 'react-forms';
import PasswordForm from 'components/utils/PasswordForm';

export default class Connections extends React.Component {

  static propTypes = {
    addTab: React.PropTypes.Func,
  }

  constructor(props) {
    super(props);
    const formValue = createValue({
      value: null,
      onChange: this.onChange.bind(this),
    });
    this.state = {
      formValue,
    };
    this.onSubmit = this.onSubmit.bind(this);
  }

  onChange(formValue) {
    this.setState({ formValue });
  }

  onSubmit = () => {
    console.log('Launching tab with: ', this.state.formValue.value);
    this.props.addTab(this.state.formValue.value);
  }
  // TODO Password fields dont seem to be playing nice with field elements.
  render() {
    return (
      <div className="customResultsSet w-30 center mt3">
        <div className="customResultsTitle">
          <div className="customTitle tc">
            Add a connection
          </div>
        </div>
        <div className="customResultsBody flex flex-row justify-center">
          <Fieldset formValue={this.state.formValue}>
            <Field select="url" label="Login URL">
              <Input className="pa1 b--none outline-transparent" />
            </Field>
            <Field select="username" label="Username" >
              <Input className="pa1 b--none outline-transparent" />
            </Field>
            <Field select="password" label="Password" Input={PasswordForm} >
              <Input className="pa1 b--none outline-transparent" />
            </Field>
            <Field select="userAgent" label="User Agent" >
              <Input className="pa1 b--none outline-transparent" />
            </Field>
            <Field select="userAgentPw" label="User Agent Password" Input={PasswordForm} >
              <Input className="pa1 b--none outline-transparent" />
            </Field>
            <Field select="version" label="Protocol Version" >
              <Input className="pa1 b--none outline-transparent" />
            </Field>
            <Field select="proxy" label="Proxy (Socks5)" >
              <Input className="pa1 b--none outline-transparent" />
            </Field>
            <Field select="id" label="Custom RETs Name" >
              <Input className="pa1 b--none outline-transparent" />
            </Field>
          </Fieldset>
        </div>
        <div className="customResultsFoot">
          <button
            onClick={this.onSubmit}
            className="customButton"
          >
            Submit
          </button>
        </div>
      </div>
    );
  }
}
