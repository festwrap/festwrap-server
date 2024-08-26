import { describe, test, expect, vi, Mock } from "vitest"
import { render } from "@testing-library/react"
import Home from "./page"
import { useSession } from "next-auth/react"

vi.mock("next-auth/react", () => {
  return {
    __esModule: true,
    useSession: vi.fn(),
  }
})

const useSessionMock = useSession as Mock<typeof useSession>

describe("Home", () => {
  test("should render sign in button when there is not session", () => {
    useSessionMock.mockReturnValue({
      update: vi.fn(),
      data: null,
      status: "unauthenticated",
    })

    const { getByText } = render(<Home />)
    expect(getByText("Not signed in")).toBeTruthy()
  })

  test("should render sign out button when there is session", () => {
    const mockSession = {
      expires: new Date(Date.now() + 2 * 86400).toISOString(),
      user: { username: "user", email: "user@gmail.com" },
    }

    useSessionMock.mockReturnValue({
      update: vi.fn(),
      data: mockSession,
      status: "authenticated",
    })

    const { getByText } = render(<Home />)

    const expectedMessage = `Signed in as ${mockSession.user.email}`
    expect(getByText(expectedMessage)).toBeTruthy()
  })
})
