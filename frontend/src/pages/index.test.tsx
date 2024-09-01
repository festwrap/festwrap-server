import { describe, test, expect, vi, afterEach } from "vitest"
import { cleanup, render, screen } from "@testing-library/react"
import Home from "."
import { useSession } from "next-auth/react"

vi.mock("next-auth/react", () => {
  return {
    __esModule: true,
    useSession: vi.fn(),
  }
})

const writeTextMock = vi.fn()

Object.assign(navigator, {
  clipboard: {
    writeText: writeTextMock,
  },
})

describe("Home", () => {
  const TOMORRROW_DATE = new Date(Date.now() + 86400).toISOString()

  afterEach(() => {
    vi.clearAllMocks()
    cleanup()
  })

  test("should render sign in button when there is not session", () => {
    vi.mocked(useSession).mockReturnValue({
      update: vi.fn(),
      data: null,
      status: "unauthenticated",
    })

    render(<Home />)

    expect(screen.getByText("Not signed in")).toBeInTheDocument()
    expect(
      screen.getByRole("button", { name: /Login with Spotify/i })
    ).toBeTruthy()
  })

  test("should render sign out button when there is session", () => {
    const mockSession = {
      expires: TOMORRROW_DATE,
      user: { name: "user", email: "user@gmail.com", accessToken: "token" },
    }

    vi.mocked(useSession).mockReturnValue({
      update: vi.fn(),
      data: mockSession,
      status: "authenticated",
    })

    render(<Home />)

    const expectedMessage = `Signed in as ${mockSession.user.email}`
    expect(screen.getByText(expectedMessage)).toBeInTheDocument()

    expect(
      screen.getByRole("button", { name: /Sign out/i })
    ).toBeInTheDocument()
  })

  test("should copy access token to clipboard", () => {
    const mockSession = {
      expires: TOMORRROW_DATE,
      user: { name: "user", email: "user@gmail.com", accessToken: "token" },
    }

    vi.mocked(useSession).mockReturnValue({
      update: vi.fn(),
      data: mockSession,
      status: "authenticated",
    })

    render(<Home />)

    const copyButton = screen.getByRole("button", {
      name: /Copy access token/i,
    })

    copyButton.click()

    expect(writeTextMock).toHaveBeenCalledWith(mockSession.user.accessToken)
  })
})
