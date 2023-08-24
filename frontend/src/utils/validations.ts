const R_PASSWORD = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/;

export type ValidationResult = {
  isValid: boolean;
  errorMessage: string;
};

function isValidEmailFormat(value: string): boolean {
  return R_PASSWORD.test(value);
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

export function validatePassword(value: string): ValidationResult {
  if (!value || value.trim() === "") {
    return {
      isValid: false,
      errorMessage: "Password is required",
    };
  }

  if (!(value.length >= 8)) {
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
