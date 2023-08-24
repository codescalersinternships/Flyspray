export type ValidationResult = {
  isValid: boolean;
  errorMessage: string;
};

function isValidEmailFormat(value: string): boolean {
  const pattern = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;
  return pattern.test(value);
}

export function validateEmail(value: string): ValidationResult {
  if (!value || value.trim() === "") {
    return {
      isValid: false,
      errorMessage: "Email is required",
    };
  }

  if (!isValidEmailFormat(value)) {
    return {
      isValid: false,
      errorMessage: "Invalid email format",
    };
  }

  return {
    isValid: true,
    errorMessage: "",
  };
}

function isValidPassword(value: string): boolean {
  return value.length >= 8;
}

function isValidPasswordRegister(value: string): boolean {
  const regexPattern =
    /^(?=.*[A-Za-z])(?=.*\d)(?=.*[_@$!%*#?&])[A-Za-z\d_@$!%*#?&]{8,}$/;
  const regexValidate: boolean = regexPattern.test(value);
  if (value.length >= 8 && regexValidate) {
    return true;
  }
  return false;
}
export function validatePassword(value: string): ValidationResult {
  if (!value || value.trim() === "") {
    return {
      isValid: false,
      errorMessage: "Password is required",
    };
  }

  if (!isValidPassword(value)) {
    return {
      isValid: false,
      errorMessage: "Password must be at least 8 characters long",
    };
  }

  return {
    isValid: true,
    errorMessage: "",
  };
}
export function validateUsername(value: string): ValidationResult {
  if (!value || value.trim() === "") {
    return {
      isValid: false,
      errorMessage: "Username is required",
    };
  }
  //TODO: add check that username does not already exist
  return {
    isValid: true,
    errorMessage: "",
  };
}

export function validatePasswordRegister(value: string): ValidationResult {
  if (!value || value.trim() === "") {
    return {
      isValid: false,
      errorMessage: "Password is required",
    };
  }

  if (!isValidPasswordRegister(value)) {
    return {
      isValid: false,
      errorMessage:
        "Password must be at least 8 characters long, contain at least one letter, one number, and one special character",
    };
  }

  return {
    isValid: true,
    errorMessage: "",
  };
}
